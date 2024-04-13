package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func installDependencies() {
	fmt.Println("Updating dependencies...")
	updateAndUpgradeCmd := exec.Command("/bin/bash", "-c", "sudo apt update && sudo apt upgrade")
	updateAndUpgradeErr := updateAndUpgradeCmd.Run()
	if updateAndUpgradeErr != nil {
		fmt.Printf("Error installing dependencies: %s\n", updateAndUpgradeErr)
	}
	fmt.Println("Dependencies updated!")

	essentials := "zsh tmux ripgrep gd fzf git build-essential"
	installEssentials := exec.Command("/bin/bash", "-c", "sudo apt install -y "+essentials)
	installEssentialsErr := installEssentials.Run()
	if installEssentialsErr != nil {
		fmt.Printf("Error installing deps: %s\n", installEssentialsErr)
	}
	fmt.Println("Installed packages: " + essentials)
}

func installFlatpakPrograms() {
	fmt.Println("Updating Flatpak...")
	applications := [8]string{
		"com.brave.Browser",
		"com.mattjakeman.ExtensionManager",
		"me.kozec.syncthingtk",
		"md.obsidian.Obsidian",
		"com.discordapp.Discord",
		"org.telegram.desktop",
		"com.getpostman.Postman",
		"it.fabiodistasio.AntaresSQL",
	}
	for _, app := range applications {
		fmt.Println("Installing " + app + " ...")
		cmd := exec.Command("flatpak install flathub " + app)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error installing: %s\n", err)
		}
		fmt.Println("Installed " + app)
	}
}

func createDirectories() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error finding home directory:\n%s", err)
	}
	paths := [3]string{
		".config", "sb", "work",
	}
	for i := range paths {
		fmt.Printf("Creating %s/%s ...\n", home, paths[i])
		err := os.MkdirAll(home+"/"+paths[i], os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s:\n%s", paths[i], err)
		}
	}
}

func installOMZ() {
	installZsh, errInstallZsh := exec.Command("/bin/bash", "-c", "sudo apt update && sudo apt install zsh && chsh -s $(which zsh) && echo Current Shell: $SHELL").Output()
	if errInstallZsh != nil {
		fmt.Printf("Error installing zsh:\n%s", errInstallZsh)
	}
	fmt.Printf("%s\n", installZsh)

	fmt.Print("Installing OhMyZsh...")
	installOhMyZsh, errInstallOhMyZsh := exec.Command("/bin/bash", "-c", `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`).Output()
	if errInstallOhMyZsh != nil {
		fmt.Printf("Error installing OhMyZsh:\n%s", errInstallOhMyZsh)
	}
	fmt.Printf("%s\n", installOhMyZsh)
}

func installNvm() {
	installCmd, installErr := exec.Command("/bin/bash", "-c", "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash").Output()
	if installErr != nil {
		fmt.Printf("Error installing NVM: %s\n", installErr)
	}
	fmt.Printf("%s\n", installCmd)

	fmt.Printf("Installing the latest LTS Node version.")

	cmd := `export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm install --lts
`
	installNode, installNodeErr := exec.Command("/bin/bash", "-c", cmd).Output()
	if installNodeErr != nil {
		fmt.Printf("%s\n", installNodeErr)
	}
	fmt.Printf("%s\n", installNode)
}

func installNeovim() {
	run := `cd ~ && curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim.appimage
chmod u+x nvim.appimage

sudo mkdir -p /opt/nvim
sudo mv nvim.appimage /opt/nvim/nvim
export PATH="$PATH:/opt/nvim/"
git clone --depth 1 https://github.com/wbthomason/packer.nvim ~/.local/share/nvim/site/pack/packer/start/packer.nvim
rm -rf ~/.config/nvim/*
git clone https://github.com/oasido/nvim.git ~/.config/nvim
`
	cmd, err := exec.Command("/bin/bash", "-c", run).Output()
	if err != nil {
		fmt.Printf("Error installing neovim:\n%s", err)
	}
	fmt.Printf("%s", cmd)

}

func main() {
	questions := []struct {
		question string
		function func()
		install  bool
	}{
		{"Install essential dependencies?", installDependencies, false},
		{"Install programs through flatpak?", installFlatpakPrograms, false},
		{"Install NVM & the latest Node LTS version?", installNvm, false},
		{"Create directories?", createDirectories, false},
		{"Install OhMyZsh?", installOMZ, false},
		{"Install Neovim?", installNeovim, false},
	}

	for i := range questions {
		confirm := askForConfirmation(questions[i].question)
		if confirm {
			questions[i].install = true
		}
	}
	for i := range questions {
		if questions[i].install {
			questions[i].function()
		}
	}
	fmt.Println("\n\nDone!\n Make sure you git clone your dotfiles, after setting up SSH.")
}
