package main

// TODO OPT: Use a library/framework for CLI apps (check out Cobra)

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

const (
	usage = `task is a CLI for managing your TODOs.

Usage:
  task [command]

Available Commands:
  add         Add a new task to your TODO list
  do          Mark a task on your TODO list as complete
  list        List all of your incomplete tasks
`
	bucket = "tasks"
)

func uitob(ui uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, ui)
	return b
}

var btoui = binary.BigEndian.Uint64

func add(db *bolt.DB, desc string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket)) // NOTE: b != nil because main makes sure the bucket exists
		id, _ := b.NextSequence()      // NOTE: safe to ignore the error (it's always nil in an Update())
		return b.Put(uitob(id), []byte(desc))
	})
}

func do(db *bolt.DB, id uint64) (string, error) {
	var desc string
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket)) // NOTE: b != nil because main makes sure the bucket exists
		idb := uitob(id)
		descb := b.Get(idb)
		if descb == nil {
			return fmt.Errorf("no task with id %d", id)
		}
		desc = string(descb) // This makes a copy of descb
		return b.Delete(idb)
	})
	if err != nil {
		return "", err
	}
	return desc, nil
}

func list(db *bolt.DB) (map[uint64]string, error) {
	tasks := make(map[uint64]string)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket)) // NOTE: b != nil because main makes sure the bucket exists
		return b.ForEach(func(idb, descb []byte) error {
			tasks[btoui(idb)] = string(descb)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println(usage)
		return
	}

	db, err := bolt.Open("tasks.db", 0600, nil)
	if err != nil {
		log.Fatalf("opening db: %v\n", err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	if err != nil {
		log.Fatalf("creating bucket: %v\n", err)
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println(usage)
			os.Exit(1)
		}
		desc := strings.Join(args[1:], " ")
		if err := add(db, desc); err != nil {
			log.Fatalf("adding task: %v\n", err)
		}
		fmt.Printf("Added %q to your task list.\n", desc)
	case "do":
		if len(args) != 2 {
			fmt.Println(usage)
			os.Exit(1)
		}
		id, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			log.Fatalf("parsing id: %v\n", err)
		}
		desc, err := do(db, id)
		if err != nil {
			log.Fatalf("completing task: %v\n", err)
		}
		fmt.Printf("You have completed the %q task.\n", desc)
	case "list":
		if len(args) > 1 {
			fmt.Println(usage)
			os.Exit(1)
		}
		tasks, err := list(db)
		if err != nil {
			log.Fatalf("listing tasks: %v\n", err)
		}
		fmt.Println("You have the following tasks:")
		for id, desc := range tasks {
			fmt.Printf("%d. %s\n", id, desc)
		}
	default:
		fmt.Println(usage)
	}
}
