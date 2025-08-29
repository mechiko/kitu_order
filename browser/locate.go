package browser

import (
	"os"
	"os/exec"
)

var defaultChromeArgs = []string{
	"--allow-insecure-localhost",
	"--disable-background-networking",
	"--disable-background-timer-throttling",
	"--disable-backgrounding-occluded-windows",
	"--disable-breakpad",
	"--disable-client-side-phishing-detection",
	"--disable-default-apps",
	"--disable-dev-shm-usage",
	"--disable-infobars",
	"--disable-extensions",
	"--disable-features=site-per-process",
	"--disable-hang-monitor",
	"--disable-ipc-flooding-protection",
	"--disable-popup-blocking",
	"--disable-prompt-on-repost",
	"--disable-renderer-backgrounding",
	"--disable-sync",
	"--disable-translate",
	"--disable-windows10-custom-titlebar",
	"--metrics-recording-only",
	"--no-first-run",
	"--no-default-browser-check",
	"--safebrowsing-disable-auto-update",
	"--enable-automation",
	"--password-store=basic",
	"--use-mock-keychain",
}
var defaultChromeArgsDisableAutomation = []string{
	"--allow-insecure-localhost",
	"--disable-background-networking",
	"--disable-background-timer-throttling",
	"--disable-backgrounding-occluded-windows",
	"--disable-breakpad",
	"--disable-client-side-phishing-detection",
	"--disable-default-apps",
	"--disable-dev-shm-usage",
	"--disable-infobars",
	"--disable-extensions",
	"--disable-features=site-per-process",
	"--disable-hang-monitor",
	"--disable-ipc-flooding-protection",
	"--disable-popup-blocking",
	"--disable-prompt-on-repost",
	"--disable-renderer-backgrounding",
	"--disable-sync",
	"--disable-translate",
	"--disable-windows10-custom-titlebar",
	"--metrics-recording-only",
	"--no-first-run",
	"--no-default-browser-check",
	"--safebrowsing-disable-auto-update",
	"--password-store=basic",
	"--use-mock-keychain",
}

// ChromeExecutable returns a string which points to the preferred Chrome
// executable file.
var ChromeExecutable = LocateChrome

// LocateChrome returns a path to the Chrome binary, or an empty string if
// Chrome installation is not found.
func LocateChrome() string {
	pathsChrome := []string{
		os.Getenv("LocalAppData") + "/Google/Chrome/Application/chrome.exe",
		os.Getenv("ProgramFiles") + "/Google/Chrome/Application/chrome.exe",
		os.Getenv("ProgramFiles(x86)") + "/Google/Chrome/Application/chrome.exe",
		os.Getenv("LocalAppData") + "/Chromium/Application/chrome.exe",
		os.Getenv("ProgramFiles") + "/Chromium/Application/chrome.exe",
		os.Getenv("ProgramFiles(x86)") + "/Chromium/Application/chrome.exe",
		os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe",
		os.Getenv("ProgramFiles") + "/Microsoft/Edge/Application/msedge.exe",
	}
	for _, path := range pathsChrome {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}

func LocateYandex() string {
	pathsChrome := []string{
		os.Getenv("LocalAppData") + "/Yandex/YandexBrowser/Application/browser.exe",
	}
	for _, path := range pathsChrome {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}

func LocateFox() string {
	pathsChrome := []string{
		os.Getenv("ProgramFiles") + "/Firefox Developer Edition/firefox.exe",
		os.Getenv("ProgramFiles(x86)") + "/Firefox Developer Edition/firefox.exe",
		os.Getenv("ProgramFiles") + "/Mozilla Firefox/firefox.exe",
		os.Getenv("ProgramFiles(x86)") + "/Mozilla Firefox/firefox.exe",
	}
	for _, path := range pathsChrome {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}

func LocateEdge() string {
	pathsChrome := []string{
		os.Getenv("ProgramFiles") + "/Microsoft/Edge/Application/msedge.exe",
		os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe",
	}
	for _, path := range pathsChrome {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}

func GetDefaultBrowser() string {
	if LocateChrome() != "" {
		return "Chrome"
	}
	if LocateEdge() != "" {
		return "MSEdge"
	}
	if LocateYandex() != "" {
		return "Yandex"
	}
	if LocateFox() != "" {
		return "Firefox"
	}
	return ""
}

func StartBrowser(browser string, url string) error {
	var err error
	// "Chrome", "Firefox", "Yandex", "MSEdge"
	switch browser {
	case "Chrome":
		args := defaultChromeArgsDisableAutomation
		appUrl := "--app=" + url
		// args = append(args, "--window-size=1000,900")
		args = append(args, appUrl)
		cmd := exec.Command(LocateChrome(), args...)
		err = cmd.Start()
	case "Firefox":
		argsFox := []string{}
		appUrlFox := url
		argsFox = append(argsFox, appUrlFox)
		cmd := exec.Command(LocateFox(), argsFox...)
		err = cmd.Start()
	case "MSEdge":
		args := defaultChromeArgs
		appUrl := "--app=" + url
		// args = append(args, "--window-size=1000,900")
		args = append(args, appUrl)
		cmd := exec.Command(LocateEdge(), args...)
		err = cmd.Start()
	case "Yandex":
		args := defaultChromeArgsDisableAutomation
		appUrl := "--app=" + url
		// args = append(args, "--window-size=1000,900")
		args = append(args, appUrl)
		cmd := exec.Command(LocateYandex(), args...)
		err = cmd.Start()
	}
	return err
}
