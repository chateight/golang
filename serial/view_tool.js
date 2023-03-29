//
// DB handler
//
'use strict';

const sqlite3 = require("sqlite3");
const db = new sqlite3.Database("./myfare.db");

// calc file name
let d = new Date();
let year = d.getFullYear();
let month = d.getMonth();
let day = d.getDate();

let date1 = new Date(year, 0, 0);
let date2 = new Date(year, month, day)
let behind = Math.floor((date2 - date1) / (24*3600*1000));      // convert milisec to day


db.serialize(() => {
    db.each("select * from tbl" + year + behind + " where stat=0", (err, row) => {
    console.log(row);
    })
});

db.close();