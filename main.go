package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	//TODO: create chip8 instance

	//TODO: Load ROM
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		//TODO: Emulator Cycle
		rl.BeginDrawing()

		//TODO: Build Texture

		rl.EndDrawing()
	}
}
