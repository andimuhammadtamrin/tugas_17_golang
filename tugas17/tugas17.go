package main

import "database/sql"
import "encoding/json"
import "fmt"
import "net/http"
import _ "mysql-master"

type daftar_buku struct{
  ID string
  Judul string
  Pengarang string
  Tahun int
}

func koneksi()(*sql.DB, error){
  db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_daftarbuku")
  if err != nil{
    return nil,err
  }
  return db,nil
}

var data []daftar_buku

func main(){
  ambil_data()
  http.HandleFunc("/daftar",ambil_daftar)
  http.HandleFunc("/cari_buku",cari_buku)
  fmt.Println("Menjalankan Web Server Pada localhost:8080")
  http.ListenAndServe(":8080",nil)
}

func ambil_daftar(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  if r.Method == "POST"{
    var result,err = json.Marshal(data)

    if err != nil{
      http.Error(w,err.Error(), http.StatusInternalServerError)
      return
    }
    w.Write(result)
    return
  }
  http.Error(w, "" , http.StatusBadRequest)
}

func cari_buku(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")

  if r.Method == "POST"{
    var namabuku = r.FormValue("Judul")
    var result []byte
    var err error

    for _, each := range data{
      if each.Judul == namabuku{
        result,err = json.Marshal(each)

        if err != nil{
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
        }
        w.Write(result)
        return
      }
    }
    http.Error(w,"Judul Buku Tidak Tersedia", http.StatusBadRequest)
    return
  }
  http.Error(w, "",http.StatusBadRequest)
}


func ambil_data(){
  db, err := koneksi()

  if err != nil{
    fmt.Println(err.Error())
    return
  }
  defer db.Close()

  rows,err := db.Query("select *from tbl_buku")
  if err != nil{
    fmt.Println(err.Error())
    return
  }
  defer rows.Close()

  for rows.Next(){
    var each = daftar_buku{}
    var err = rows.Scan(&each.ID,&each.Judul,&each.Pengarang,&each.Tahun)

    if err != nil{
      fmt.Println(err.Error())
      return
    }
    data = append(data,each)
  }
  if err = rows.Err(); err != nil{
    fmt.Println(err.Error())
    return
  }

}
