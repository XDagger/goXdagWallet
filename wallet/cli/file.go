package cli

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	xdagoUtils "goXdagWallet/xdago/utils"
)

var validate = func(input string) error {
	if !xdagoUtils.FileExists(input) {
		return errors.New("input path not exists")
	}
	return nil
}

func inputFilePath() {
	prompt := promptui.Prompt{
		Label:    "Path to Mnemonic text file",
		Validate: validate,
	}

	result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input path failed %v\n", err)
		result, err = prompt.Run()
	}
	fmt.Println(result)

}

func inputFolderPath() {
	prompt := promptui.Prompt{
		Label:    "Path to Non Mnemonic data folder",
		Validate: validate,
	}

	result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input path failed %v\n", err)
		result, err = prompt.Run()
	}
	fmt.Println(result)
}
