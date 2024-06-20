package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	chip8 "github.com/nhlmg93/chimp_ate/pkg/chip8"
)

func loadROM(chip *chip8.Chip8) {
	projDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory", err)
		return
	}
	romPath := filepath.Join(projDir, "chip8-test-rom", "test_opcode.ch8")
	rom, err := os.Open(romPath)
	if err != nil {
		fmt.Println("Error Opening file", err)
		return
	}
	defer rom.Close()

	content, err := io.ReadAll(rom)
	if err != nil {
		fmt.Println("Error Reading file", err)
		return
	}
	for i := 0; i < len(content); i++ {
		chip.Memory[i+0x200] = content[i]
	}

}

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	chip := chip8.NewChip8()
	loadROM(chip)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		chip.Cycle()
		rl.BeginDrawing()
		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				if chip.Graphics[y*64+x] == 1 {
					rl.DrawRectangle(int32(x)*10, int32(y)*10, 10, 10, rl.RayWhite)
				} else {
					rl.DrawRectangle(int32(x)*10, int32(y)*10, 10, 10, rl.Black)
				}
			}

		}
		rl.EndDrawing()
	}
}
