var express = require('express');
var router = express.Router();

let nameArray = new Array();

/* GET home page. */
router.get('/', function(req, res, next) {

  let cTime = new Date();
  let update = "setTimeout(function () {location.reload();}, 2000)"

async function dbReadExec(){
  try{
    nameArray = await dbRead();
    
    res.render('index', { title: 'Welcome to CoderDojo', date: cTime, db: nameArray, update: update});
  } catch (err){
    console.log('error');
  } finally {
    console.log("done");
  }
}

dbReadExec();

});

//
// DB handler
//
function dbRead(){
  return new Promise((resolve, reject) => {

  const sqlite3 = require("sqlite3");
  const db = new sqlite3.Database("/Users/usamiryuuichi/go/src/serial/myfare.db");
  
  // calc file name
  let d = new Date();
  let year = d.getFullYear();
  let month = d.getMonth();
  let day = d.getDate();
  
  let date1 = new Date(year, 0, 0);
  let date2 = new Date(year, month, day)
  let behind = Math.floor((date2 - date1) / (24*3600*1000));      // convert milisec to day
  
  nameArray = [];       // clear Array
  
  db.serialize(() => {

  db.all("select * from tbl" + year + behind + " where stat=1 order by time desc", (err, rows) =>{
    rows.forEach(row => nameArray.push(row.name));
    //console.log(nameArray);
    resolve(nameArray);
})

  db.close();

  });
  });
}

module.exports = router;
