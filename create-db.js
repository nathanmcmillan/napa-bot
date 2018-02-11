
const fs = require('fs');
const sql = require('sqlite3');
const file = fs.readFileSync('db.sql').toString();
const name = 'trade.db';

console.log('creating database');

const db = new sql.Database(name, (error) => {
    if (error) {
        console.error(error.message);
        return;
    }
    console.log('connected');
});

db.exec(file, (error) => {
    if (error) {
        console.error(error.message);
        return;
    }
    console.log('executed');
});

db.close((error) => {
    if (error) {
        console.error(error.message);
        return;
    }
    console.log('closed');
});
