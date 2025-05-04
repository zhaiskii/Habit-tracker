package storage

import(
	"fmt"
	"sync"
	"net/http"
	"database/sql"
	"habit21/internal/config"
	_ "github.com/lib/pq"
)

type Storage struct{
	DB *sql.DB
	mu sync.Mutex
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "internal/storage/New"
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name)
	conn, err := sql.Open("postgres", psql)
	if err!=nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &Storage{DB: conn}, nil
}

func (conn *Storage) Create(id int, name string, progress int, completed bool) error {
	conn.mu.Lock()
	fmt.Println(name)
	defer conn.mu.Unlock()
	_, err := conn.DB.Exec("INSERT INTO habit (id, habit, progress, completed) VALUES ($1, $2, $3, $4)", id, name, progress, completed)
	if err!=nil {
		return err
	}
	
	return err
}

func (conn *Storage) Delete(id string) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	res, err := conn.DB.Exec("DELETE FROM habit WHERE id=$1", id)
	if fi,_:=res.RowsAffected(); fi==0 {
		return fmt.Errorf("Not Found", http.StatusNotFound)
	}
	return err
}

func (conn *Storage) Show() (*sql.Rows, error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	rows, err := conn.DB.Query("SELECT * FROM habit")
	if err!=nil {
		return nil, err
	}
	return rows, nil
}

func (conn *Storage) Update(id string, progress int, completed bool, date string) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	_, err := conn.DB.Exec("UPDATE habit SET progress=$1, completed=$2 WHERE id=$3", progress, completed, id)
	var cnt int
	rows, err := conn.DB.Query("SELECT count FROM days WHERE date=$1", date)
	rows.Next()
	err = rows.Scan(&cnt)
	fmt.Println(date, cnt, err)
	if err!=nil {
		cnt=0
		fmt.Println("whyyy!!!/???")
		_, err := conn.DB.Exec("INSERT INTO days (date, count) VALUES($1, $2)", date, cnt)
		if err!=nil {
			return err
		}
	}
	if completed==true {
		cnt++
	} else {
		cnt--
	}
	_, err = conn.DB.Exec("UPDATE days SET count = $1 WHERE date=$2", cnt, date)
	return err
} 

func (conn *Storage) UpdateDefault() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	_, err := conn.DB.Exec("UPDATE habit SET completed=$1", false)
	return err
} 

func (conn *Storage) ShowTable() (*sql.Rows, error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	rows, err := conn.DB.Query("SELECT * FROM days")
	if err!=nil {
		return nil, err
	}
	return rows, nil
}

func (conn *Storage) AddUser(email string, hashed string) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	_, err := conn.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, hashed)
	return err
} 