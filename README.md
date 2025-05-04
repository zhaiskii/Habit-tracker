Habit Tracker

An app for tracking habits

Prerequisites:
1) Install golang
2) Install psql (or other tool for database interaction) and postgres for working with database

Setting up:
1) clone the repo: git clone https://github.com/zhaiskii/Habit-tracker.git
2) Create a database in your local machine, the needed information related to the name and port of database is provided at the yaml file. If needed you can change it
3) Create 3 tables at your database:
   
   CREATE TABLE days(date TEXT,count INT);

   CREATE TABLE days(id INT PRIMARY KEY, habit TEXT, progress INT, completed BOOL);

   CREATE TABLE users(email TEXT,password TEXT);

5) Run the code by the command go run ./cmd/main.go , to open the frontend on browser click on go live in VScode (if you are not in vscode there should be way to open frontend in some port, but don't open frontend just from file)
6) Make sure the port 8080 is free

Tried to separate my code in different files to make component-based design. Used chi router and CORS to permit the request to the server from another origin. Haven't finished authentication.

The main drawback is that login/registration don't work, you can even not test it, it doesn't work. Other from that everything is good with the app.

Chosen this stack because I am familiar to GO and love coding in it, and I am a complete newbie in frontend and therefore chose the easiest stack. 
