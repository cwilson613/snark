package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/google/go-github/v43/github"
	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

type Release struct {
	Name      string
	ID        int64
	AssetsURL string
}

var Version string = ""
var BuildTime string = ""
var access_token string

func main() {
	argNum := len(os.Args)
	if argNum > 1 {
		arg1 := os.Args[1]

		switch {
		//init
		case arg1 == "list":
			// only bother with listing GA releases by default
			ga := true
			if argNum == 3 {
				if os.Args[2] != "all" {
					log.Fatal("Unrecognized flag: " + os.Args[2])
				} else {
					ga = false
				}
			}
			fmt.Printf("Available DKP Versions:\n")
			listDKPVersions(ga, "")
		//bootstrap
		case arg1 == "init":
			version := checkVersion()
			if version == "" {
				fmt.Println("No previous DKP use detected, please select desired version: ")
				listDKPVersions(true, "")
			} else {
				listDKPVersions(true, version)
			}
		case arg1 == "version":
			fmt.Println("Version:	" + Version)
			fmt.Println("		" + BuildTime)
			fmt.Println("OS:		" + runtime.GOOS)
			fmt.Println("Arch:		" + runtime.GOARCH)
		case arg1 == "help":
			printDefaultUsage()
		default:
			if argNum == 1 {
				log.Fatal("No DKP Version Provided")
				printDefaultUsage()
			} else if argNum == 2 {
				version := os.Args[1]
				listDKPVersions(true, version)
			}
		}

	} else {
		printDefaultUsage()
	}
}

func checkVersion() string {
	// Try to read local .dkp file and parse version from it
	version := ""
	content, err := ioutil.ReadFile(".dkp")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return version
		}
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	return string(content)
}

func printDefaultUsage() {
	fmt.Printf("Usage:\n" +
		" dkpswitch 				prints usage\n" +
		" dkpswitch version			prints version information\n" +
		" dkpswitch list 			list available GA DKP versions\n" +
		" dkpswitch list all		list all DKP versions\n" +
		" dkpswitch <version>		switch to provided GA DKP version [e.g. 'v2.1.1']\n" +
		" dkpswitch init			switch to DKP version previously used in this directory.")
}

func listDKPVersions(ga bool, specificVersion string) {
	rs := githubAuth()

	dkpTags := []github.RepositoryRelease{}
	/*
		So as it turns out, many of the early Konvoy v1.x releases are compressed using
		Bzip2, which - while effective - is a completely different logical path from
		standard Gzip compression. For now let's focus on the present 2.x line, and
		if we ever need to include this functionality, it's already stubbed here.
	*/
	dkpTags = append(dkpTags, getDKP1tags(rs, ga)...)
	dkpTags = append(dkpTags, getDKP2tags(rs, ga)...)

	names := []string{}
	for _, tag := range dkpTags {
		names = append(names, *tag.Name)
	}

	// Silly way to reverse-sort a slice of strings...
	// but easier to read than some looping nonsense
	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	if specificVersion == "" {
		prompt := promptui.Select{
			Label: "Select DKP Version",
			Items: names,
		}

		_, result, err := prompt.Run()

		if err != nil {
			log.Fatal(err)
		}

		for _, tag := range dkpTags {
			if *tag.Name == result {
				setupDKP(tag)
			}
		}
	} else {
		found := false
		for _, tag := range dkpTags {
			if strings.Contains(*tag.Name, specificVersion) {
				found = true
				setupDKP(tag)
			}
		}
		if !found {
			log.Fatal("Could not find a valid GA release for " + specificVersion)
		}
	}
}

func readGithubAccessToken() string {
	token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	// fmt.Println("GITHUB_PERSONAL_ACCESS_TOKEN: " + token)
	if token == "" {
		byteToken, _ := terminal.ReadPassword(int(syscall.Stdin))
		token = string(byteToken)
		fmt.Println("\nTo avoid having to input your access token each run, consider setting the GITHUB_PERSONAL_ACCESS_TOKEN environment variable ")
	}
	return token
}

func githubAuth() github.RepositoriesService {
	token := readGithubAccessToken()
	access_token = token

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	// useful for debugging auth
	fmt.Println("Welcome " + *user.Login + "!")

	return *client.Repositories

}

func getDKP1tags(rs github.RepositoriesService, ga bool) []github.RepositoryRelease {
	return getTags(rs, "konvoy", ga)
}
func getDKP2tags(rs github.RepositoriesService, ga bool) []github.RepositoryRelease {
	return getTags(rs, "konvoy2", ga)
}

func getTags(rs github.RepositoriesService, repoName string, ga bool) []github.RepositoryRelease {
	releaseList, _, err := rs.ListReleases(context.Background(), "mesosphere", repoName, nil)

	validReleases := []github.RepositoryRelease{}

	if err != nil {
		log.Fatal(err)
		return validReleases
	}

	for _, s := range releaseList {
		if ga {
			if !strings.Contains(*s.Name, "-") {
				validReleases = append(validReleases, *s)
			}
		} else {
			validReleases = append(validReleases, *s)
		}
	}

	return validReleases
}

func checkCurrentVersions() {

	// DKP 1.x
	_, kerr := exec.LookPath("konvoy")
	if kerr != nil {
		fmt.Println("Could not find konvoy binary in $PATH: " + kerr.Error())
	} else {
		cmd := exec.Command("konvoy version")
		stdout, err := cmd.Output()
		if err != nil {
			log.Fatal("Could not get konvoy version: " + err.Error())
		} else {
			fmt.Println(stdout)
		}
	}

	// DKP 2.x
	_, dkperr := exec.LookPath("dkp")
	if dkperr != nil {
		fmt.Println("Could not find konvoy binary in $PATH: " + dkperr.Error())
	} else {
		cmd := exec.Command("dkp version")
		stdout, err := cmd.Output()
		if err != nil {
			log.Fatal("Could not get DKP version: " + err.Error())
		} else {
			fmt.Println(stdout)
		}
	}
}

func getSpecificRelease(version string) {

}

func setupDKP(release github.RepositoryRelease) {
	// append the v if missing
	if !strings.HasPrefix(*release.Name, "v") {
		*release.Name = "v" + *release.Name
	}

	checkTmpDir()
	// checkBinaryLinks(release)
	binPath, linkPath := downloadRelease(release)

	// fmt.Println("Creating link: " + linkPath + " -> " + binPath)
	if _, err := os.Lstat(linkPath); err == nil {
		os.Remove(linkPath)
	}
	err := os.Symlink(binPath, linkPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Version set to " + *release.Name + "!")

	// write the version set at ./.dkp
	_ = os.Remove(".dkp")
	d1 := []byte(*release.Name)
	err = os.WriteFile(".dkp", d1, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func checkTmpDir() {
	nerr := os.MkdirAll("/tmp/dkp", 0750)
	if nerr != nil {
		log.Fatal(nerr)
	}
}

func downloadRelease(release github.RepositoryRelease) (string, string) {

	osName := runtime.GOOS
	var url string
	var filename string
	var linkPath string
	var binaryPath string
	skipDownload := false
	skipInflation := false

	for _, asset := range release.Assets {
		if strings.Contains(*asset.BrowserDownloadURL, osName) {
			url = *asset.URL
			filename = path.Base(*asset.BrowserDownloadURL)
		}
	}

	// fmt.Println("Download URL: " + url)

	// strip out the prefix URL garbage
	filepath := "/tmp/dkp/" + filename

	err := os.MkdirAll("/tmp/dkp/"+*release.Name, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// have we already downloaded this?
	_, err = os.Stat(filepath)
	if err == nil {
		skipDownload = true
	}

	// have we already inflated this?
	_, err = os.Stat("/tmp/dkp/" + *release.Name)
	if err == nil {
		if strings.Contains(filepath, "konvoy") {
			_, kerr := os.Stat("/tmp/dkp/" + *release.Name + "konvoy")
			_, perr := os.Stat("/tmp/dkp/" + *release.Name + "konvoy-preflight")
			if kerr == nil && perr == nil { // konvoy binaries exist
				skipInflation = true
			}
		} else {
			_, derr := os.Stat("/tmp/dkp/" + *release.Name + "dkp")
			if derr == nil { // dkp binary exists
				skipInflation = true
			}
		}
	}

	//set the proper binary and link names in advance
	if strings.Contains(filepath, "konvoy") {
		linkPath = "/usr/local/bin/konvoy"
		binaryPath = "/tmp/dkp/" + *release.Name + "/konvoy"
	} else {
		linkPath = "/usr/local/bin/dkp"
		binaryPath = "/tmp/dkp/" + *release.Name + "/dkp"
	}

	if skipDownload {
		fmt.Println(*release.Name + " is already in /tmp/dkp - skipping download.")
	} else {

		out, err := os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// Create new httpclient to fetch the release file
		client := http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		// set proper headers (auth, content type, user-agent request)
		req.Header = http.Header{
			"Accept": []string{
				"application/octet-stream",
			},
			"User-Agent":    []string{"request module"},
			"Authorization": []string{"token " + access_token},
		}

		// execute
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// write the file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	if skipInflation {
		fmt.Println("Skipping inflation, already exists!")
	} else {
		// fmt.Println("tar -xvf " + filepath)
		// fmt.Println("/tmp/dkp/" + *release.Name)
		cmd := exec.Command("tar", "-xvf", filepath)
		cmd.Dir = "/tmp/dkp/" + *release.Name
		_, err := cmd.Output()

		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(stdout)

		if strings.Contains(filepath, "konvoy") {
			fmt.Println("Extra konvoy steps")
			for _, bin := range []string{"konvoy", "konvoy-preflight"} {
				copyFrom := "/tmp/dkp/" + *release.Name + "/konvoy_" + *release.Name + "/" + bin
				copyTo := "/tmp/dkp/" + *release.Name + "/" + bin
				fmt.Println("mv " + copyFrom + " " + copyTo)
				err = os.Rename(copyFrom, copyTo)
				if err != nil {
					log.Fatal("Failed to copy old style konvoy binary to proper release dir: " + err.Error())
				}
			}
			// now get rid of the spurious internal folder
			fmt.Println("rm /tmp/dkp/" + *release.Name + "/konvoy_" + *release.Name)
			err = os.RemoveAll("/tmp/dkp/" + *release.Name + "/konvoy_" + *release.Name)
			if err != nil {
				log.Fatal(err)
			}

			// get rid of the download
			// NOTE: in future want to keep it and simply reinflate if folder isn't there but the tarball IS
			fmt.Println("rm " + filepath)
			err = os.Remove(filepath)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	return binaryPath, linkPath

}
