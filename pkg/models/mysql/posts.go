package mysql

import (
	"database/sql"
	"errors"

	"github.com/DTGlov/goweb.git/pkg/models"
)

//Define a PostModel type which wraps a sql.DB connection pool.
type PostModel struct {
	DB *sql.DB
}

//This will insert  a new post into the database
func (p *PostModel) Insert(title, content, expires string) (int, error) {
	//query statement to insert into the database
	stmt := `INSERT INTO posts (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//use the Exec() method on the embedded connection pool to execute the statement
	result, err := p.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}

	//Use the lastId method to get the ID of the last inserted post
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	//the id returned is an int64 so we convert into an int type
	return int(id), nil
}

func (p *PostModel) Get(id int) (*models.Post, error) {
	//Query to get a specific post
	stmt := `SELECT id, title, content, created, expires FROM posts
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := p.DB.QueryRow(stmt, id)

	//intialize a pointer to a new zeroed Post struct
	c := &models.Post{}

	//use rows.Scan() to copy each field in the row to their corresponding fields in the struct
	err := row.Scan(&c.ID, &c.Title, &c.Content, &c.Created, &c.Expires)
	if err != nil {
		//if query returns no rows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	//if everything went OK then return the Post object
	return c, nil
}

func (p *PostModel) Latest() ([]*models.Post, error) {
	//query to get multiple rows
	stmt := `SELECT id,title,content,created,expires FROM posts WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	//use the Query method on the db pool to return multiple rows
	rows, err := p.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	//we defer.Close() to ensure the sql.Rows resultset is always properly closed before the Latest() method returns.
	defer rows.Close()

	//Initialize an empty slice to hold the models.Post objects
	posts := []*models.Post{}

	//use rows.Next to scan through the resultset.
	for rows.Next() {
		//create a pointer to the new  zeroed  Post struct
		s := &models.Post{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	//after the scan is over we call rows.Err() to retrieve the error that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	//if everything is ok we
	return posts, nil

}

// func PrintMost() {
// 	var mostPosts []*models.Post

// 	posts := &models.Post{
// 		ID:      1,
// 		Title:   "Yes",
// 		Content: "Hujj",
// 	}

// 	mostPosts = append(mostPosts, posts)

// 	for _, v := range mostPosts {
// 		fmt.Println(*v)
// 	}
// }
