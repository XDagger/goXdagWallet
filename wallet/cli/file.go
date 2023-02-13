package cli

import (
	"errors"
	"fmt"
	xdagoUtils "goXdagWallet/xdago/utils"

	"github.com/manifoldco/promptui"
)

var validate = func(input string) error {
	if !xdagoUtils.FileExists(input) {
		return errors.New("input path not exists")
	}
	return nil
}

func inputFilePath() string {
	prompt := promptui.Prompt{
		Label:    "Path to Mnemonic text file",
		Validate: validate,
	}

	result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input path failed %v\n", err)
		result, err = prompt.Run()
	}
	return result

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
