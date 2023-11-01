package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	dbFileName = "tasks.db"
)

var tasksBucket = []byte("tasks")

// json annotations, field tags
type Task struct {
	ID        int    `json: "id"`
	Details   string `json: "details"`
	Completed bool   `json: "completed"`
}

func CreateTask(task *Task) error {
	return withTasksDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			// Generate ID for the user.
			// This returns an error only if the Tx is closed or not writeable.
			// That can't happen in an Update() call so I ignore the error check.
			bucket := tx.Bucket(tasksBucket)
			id, err := bucket.NextSequence()
			if err != nil {
				return err
			}

			task.ID = int(id)

			// Marshal user data into bytes.
			buf, err := json.Marshal(&task)
			if err != nil {
				return err
			}

			// Persist bytes to taks bucket.
			return bucket.Put(itob(task.ID), buf)
		})
	})
}

func ListTasks(completed bool) ([]*Task, error) {
	var tasks []*Task
	return tasks, withTasksDB(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(tasksBucket)

			return bucket.ForEach(func(k, b []byte) error {
				var task Task
				if err := json.Unmarshal(b, &task); err != nil {
					return err
				}
				// skip completed
				if task.Completed != completed {
					return nil
				}
				tasks = append(tasks, &task)
				return nil
			})
		})
	})
}

func MarkTaskAsCompleted(task *Task) error {
	return withTasksDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(tasksBucket)

			b := bucket.Get(itob(task.ID))
			if b == nil {
				return fmt.Errorf("task not found with Id=%v", task.ID)
			}
			// deserialize
			if err := json.Unmarshal(b, task); err != nil {
				return err
			}

			task.Completed = true
			// serialize
			b, err := json.Marshal(&task)

			if err != nil {
				return err
			}

			return bucket.Put(itob(task.ID), b)
		})
	})
}

func DeleteTask(task *Task) error {
	return withTasksDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(tasksBucket)

			b := bucket.Get(itob(task.ID))
			if b == nil {
				return fmt.Errorf("task not found with Id=%v", task.ID)
			}

			if err := json.Unmarshal(b, task); err != nil {
				return err
			}

			return bucket.Delete(itob(task.ID))
		})
	})
}

func withTasksDB(fun func(db *bolt.DB) error) error {
	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		return err
	}

	defer db.Close()

	// every time make a quick check to see if the bucket exists
	// Retrieve the tasks bucket.
	// This should be created when the DB is first opened.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(tasksBucket)
		return err
	})

	if err != nil {
		return err
	}

	return fun(db)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
