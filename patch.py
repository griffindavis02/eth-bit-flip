import os
import sys
import re

dirGeth = re.sub(r'\\', '/', sys.argv[1])
fileTriggers = {
    'cmd/utils/flags.go': ['pcsclite "github.com/gballet/go-libpcsclite"', 'Usage: "Catalyst mode (eth2 integration testing)",'],
    'cmd/geth/main.go': ['utils.MetricsInfluxDBOrganizationFlag,', 'app.Flags = append(app.Flags, metricsFlags...)']
}

patches = {
    'cmd/utils/flags.go': ['bitflip "github.com/griffindavis02/eth-bit-flip/flags"',
    """
    // Flags for simulating soft errors in the blockchain
	FlipInitialized = bitflip.FlipInitialized

	FlipPath = bitflip.FlipPath

	FlipStart = bitflip.FlipStart

	FlipStop = bitflip.FlipStop

	FlipRestart = bitflip.FlipRestart

	FlipType = bitflip.FlipType

	FlipCounter = bitflip.FlipCounter

	FlipIterations = bitflip.FlipIterations

	FlipVariables = bitflip.FlipVariables

	FlipDuration = bitflip.FlipDuration

	FlipTime = bitflip.FlipTime

	FlipRate = bitflip.FlipRate

	FlipRates = bitflip.FlipRates

	FlipPost = bitflip.FlipPost

	FlipHost = bitflip.FlipHost
    """],

    'cmd/geth/main.go': ["""
    flipFlags = []cli.Flag{
		utils.FlipPath,
		utils.FlipStart,
		utils.FlipStop,
		utils.FlipRestart,
		utils.FlipType,
		utils.FlipCounter,
		utils.FlipIterations,
		utils.FlipVariables,
		utils.FlipDuration,
		utils.FlipTime,
		utils.FlipRate,
		utils.FlipRates,
		utils.FlipPost,
		utils.FlipHost,
	}
    """,
    '\tapp.Flags = append(app.Flags, flipFlags...)']
}

def main():
    """
    Parse through the go-ethereum source code and insert the necessary function
    calls
    """
    fileNames = list(fileTriggers)

    triggerContent = findTrigger(fileNames[0], fileTriggers[fileNames[0]][0], 0)
    patch(fileNames[0], triggerContent[1], patches[fileNames[0]][0], triggerContent[0], 1)
    triggerContent = findTrigger(fileNames[0], fileTriggers[fileNames[0]][1], 0)
    patch(fileNames[0], triggerContent[1], patches[fileNames[0]][1], triggerContent[0], 0)

def findTrigger(fileName: str, trigger: str, overrides: int) -> (int, list[str]):
    """
    Reads the selected file in the geth soruce code and returns the line number
    of the selected trigger and a list of file lines.
    """
    filePath = os.path.join(dirGeth, fileName)
    print(f'Reading file {filePath}')

    with open(filePath, 'r') as rFile:
        lines = rFile.readlines()
        overrideCheck = 0
        for lineNum, line in enumerate(lines):
            if trigger in line:
                if overrideCheck >= overrides:
                    print('found trigger')
                    return (lineNum, lines)
                overrideCheck+=1

def patch(fileName: str, fileContents: list[str], patchStr: str,
            lineNum: int, offset: int) -> None:
    """
    Takes the given file name and its contents, writes the patch into the file
    at line lineNum + offset
    """
    filePath = os.path.join(dirGeth, fileName)
    print(f'Writing file {filePath}')
    with open(filePath, 'w') as wFile:
        fileContents[lineNum+offset+1:lineNum+offset+1] = re.split(r'(\n)', patchStr)
        wFile.write(''.join(fileContents))

main()