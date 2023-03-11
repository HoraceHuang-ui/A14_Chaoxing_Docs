package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

// Exam
func main() {
	db, err := sql.Open("mysql", "root:lyh701721@tcp(47.99.66.125:5721)/fwwb")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	fmt.Println("Connection established")

	// Retrieve exam-clazz relations
	clazzes_table, err := db.Query("SELECT exam_id, clazzid FROM t_exam_relation")
	if err != nil {
		fmt.Println(err)
	}
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
	clazzes_table.Close()

	fmt.Println("Retrieve relations complete\n")

	for clazz, exams := range clazzes {
		f := excelize.NewFile()
		f.SetSheetRow("Sheet1", "A1", &[]interface{}{"personid", "姓名", "班级", "课程id", "课程"})
		col := 'F'
		for _, e := range exams {
			f.SetCellValue("Sheet1", fmt.Sprintf("%c1", col), e)
			col++
		}

		fmt.Printf("%s Set column title complete\n", clazz)

		clazz_stus, err := db.Query("SELECT personid, courseid FROM t_person_course_clazz_role WHERE clazzid=?", clazz)
		if err != nil {
			fmt.Println("clazz_stus query: ", err)
		}
		defer clazz_stus.Close()
		row := 2
		for clazz_stus.Next() {
			// Fill columns A,B,C,D,E
			id := ""
			name := ""
			coursename := ""
			courseid := ""
			clazz_stus.Scan(&id, &courseid)
			course_relations, _ := db.Query("SELECT name FROM t_course WHERE courseid=?", courseid)
			course_relations.Next()
			course_relations.Scan(&coursename)
			course_relations.Close()
			stuinfo, err := db.Query("SELECT user_name FROM t_person WHERE personid=?", id)
			if err != nil {
				fmt.Println("stuinfo query: ", err)
			}
			stuinfo.Next()
			err = stuinfo.Scan(&name)
			if err != nil {
				fmt.Println("stuinfo scan: ", err)
			}
			stuinfo.Close()
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), id)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), name)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), clazz)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), courseid)
			f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), coursename)

			// Retrieve scores for each exam
			col := 'F'
			for _, e := range exams {
				var score float64 = 0.0
				scoredb, err := db.Query("SELECT score FROM t_exam_answer WHERE personid=? AND exam_id=?", id, e)
				if err != nil {
					fmt.Println(err)
				} else {
					scoredb.Next()
					scoredb.Scan(&score)
					f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", col, row), score)
				}
				scoredb.Close()
				col++
			}

			fmt.Printf("%s Fill value complete \n", name)

			row++
		}
		if err := f.SaveAs(fmt.Sprintf("exams\\exam_%s.xlsx", clazz)); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s Save sheet complete\n\n", clazz)
	}
}

/*
// Homework
func main() {
	db, err := sql.Open("mysql", "root:lyh701721@tcp(47.99.66.125:5721)/fwwb")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	fmt.Println("Connection established")

	// Retrieve work-clazz relations
	clazzes_table, err := db.Query("SELECT work_id, clazzid FROM t_work_relation")
	if err != nil {
		fmt.Println(err)
	}
	clazzes := map[string][]string{}
	for clazzes_table.Next() {
		c, w := "", ""
		clazzes_table.Scan(&w, &c)
		if _, ok := clazzes[c]; !ok {
			clazzes[c] = []string{w}
		} else {
			tmp := clazzes[c]
			clazzes[c] = append(tmp, w)
		}
	}
	clazzes_table.Close()

	fmt.Println("Retrieve relations complete\n")

	for clazz, works := range clazzes {
		f := excelize.NewFile()
		f.SetSheetRow("Sheet1", "A1", &[]interface{}{"personid", "姓名", "班级", "课程id", "课程"})
		col := 'F'
		for _, w := range works {
			f.SetCellValue("Sheet1", fmt.Sprintf("%c1", col), w)
			col++
		}

		fmt.Printf("%s Set column title complete\n", clazz)

		clazz_stus, err := db.Query("SELECT personid, courseid FROM t_person_course_clazz_role WHERE clazzid=?", clazz)
		if err != nil {
			fmt.Println("clazz_stus query: ", err)
		}
		defer clazz_stus.Close()
		row := 2
		for clazz_stus.Next() {
			// Fill columns A,B,C,D,E
			id := ""
			name := ""
			coursename := ""
			courseid := ""
			clazz_stus.Scan(&id, &courseid)
			course_relations, _ := db.Query("SELECT name FROM t_course WHERE courseid=?", courseid)
			course_relations.Next()
			course_relations.Scan(&coursename)
			course_relations.Close()
			stuinfo, err := db.Query("SELECT user_name FROM t_person WHERE personid=?", id)
			if err != nil {
				fmt.Println("stuinfo query: ", err)
			}
			stuinfo.Next()
			err = stuinfo.Scan(&name)
			if err != nil {
				fmt.Println("stuinfo scan: ", err)
			}
			stuinfo.Close()
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), id)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), name)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), clazz)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), courseid)
			f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), coursename)

			// Retrieve scores for each homework
			col := 'F'
			for _, w := range works {
				var score float64 = 0.0
				scoredb, err := db.Query("SELECT score FROM t_work_answer WHERE personid=? AND work_id=?", id, w)
				if err != nil {
					fmt.Println(err)
				} else {
					scoredb.Next()
					scoredb.Scan(&score)
					f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", col, row), score)
				}
				scoredb.Close()
				col++
			}

			fmt.Printf("%s Fill value complete \n", name)

			row++
		}
		if err := f.SaveAs(fmt.Sprintf("works\\work_%s.xlsx", clazz)); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s Save sheet complete\n\n", clazz)
	}
}
*/
