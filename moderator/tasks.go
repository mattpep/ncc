package moderator

import "ncc/db"
import (
	"fmt"
	"log"

	"github.com/mattn/go-tty"
)

func ShowTasks() {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	tasks, err := db.GetModerationTasks()
	if err != nil {
		fmt.Printf("Could not get moderation tasks from db: %v\n", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("Nothing to do")
		return
	}

	for _, task := range tasks {
		fmt.Printf("Comment in need of approval: \n PostRef: %s\n Author:  %s at %s\n From:    %s\n Comment: %s\n", task.PostRef, task.DisplayName, task.DateTime, task.Actor, task.Body)
		fmt.Printf("What action? ('a' == Approve, 'd' to delete, 'q' to quit, any other key to skip)  ")
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		var choice = string(r)
		choice = string(r)

		if choice == "a" || choice == "A" {
			err = db.ApproveComment(task.Id)
			if err != nil {
				fmt.Printf("Could not approve: %v", err)
				return
			}
			fmt.Println("Approved")
		} else if choice == "d" || choice == "D" {
			err = db.DeleteComment(task.Id)
			if err != nil {
				fmt.Printf("Could not delete: %v", err)
				return
			}
			fmt.Println("Deleted")
		} else if choice == "q" || choice == "Q" {
			fmt.Println("\nStopping")
			return
		} else {
			fmt.Println("Skipping")
		}
		fmt.Printf("\n\n")
	}
	fmt.Println("All done!")
}
