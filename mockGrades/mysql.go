package mysql

import (
    "fmt"
    "database/sql"

    // "github.com/xuri/excelize/v2"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:password@tcp()/fwwb")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()

    // Retrieve exam-clazz relations
    clazzes_table, _ := db.Query("SELECT exam_id, clazz_id FROM t_exam_relation")
    clazzes := map[string][]string{}
    for clazzes_table.Next() {
        c, e := "", ""
        clazzes_table.Scan(&e, &c)
        if _, ok := clazzes[c]; !ok {
            clazzes[c] = []string{e}
        } else {
            tmp := clazzes[c]
            clazzes[c] = append(tmp, e)
        }
    }

    for clazz, exams := range clazzes {
        f := excelize.NewFile()
        f.SetSheetRow("Sheet1", "A1", &[]interface{"personid", "姓名", "班级"})
        col := 'D'
        for _, e := range exams {
            f.SetCellValue("Sheet1", fmt.Sprintf("%c1",col), e)
            col++
        }

        /*
        // fill columns A,B,C
        clazz_stus, _ := db.Query("SELECT personid FROM t_person_course_clazz_role WHERE clazzid=?", clazz)
        row := 2
        for clazz_stus.Next() {
            id := ""
            clazz_stus.Scan(&id)
            stuinfo, _ := db.Query("SELECT username FROM t_person WHERE personid=?", id)
            name := ""
            stuinfo.Scan(&name)
            f.SetCellValue("Sheet1", fmt.Sprintf("A%d",row), id)
            f.SetCellValue("Sheet1", fmt.Sprintf("B%d",row), name)
            f.SetCellValue("Sheet1", fmt.Sprintf("C%d",row), clazz)
            row++
        }
        */

        for _, e := range exams {
            exam_answers, _ := db.Query("SELECT personid, score FROM t_exam_answer WHERE exam_id=?", e)
            
        }

    }
}