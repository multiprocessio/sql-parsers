package main

import (
	"fmt"
	
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/parse"
	"github.com/kr/pretty"
)

func main() {
	simple := "SELECT * FROM x"
	
	medium := "SELECT COUNT(1) AS count, name section FROM t GROUP BY name LIMIT 10"
	
	complex := `
SELECT 
	country.country_name_eng,
	SUM(CASE WHEN call.id IS NOT NULL THEN 1 ELSE 0 END) AS calls,
	AVG(ISNULL(DATEDIFF(SECOND, call.start_time, call.end_time),0)) AS avg_difference
FROM country 
LEFT JOIN city ON city.country_id = country.id
LEFT JOIN customer ON city.id = customer.city_id
LEFT JOIN call ON call.customer_id = customer.id
GROUP BY 
	country.id,
	country.country_name_eng
HAVING AVG(ISNULL(DATEDIFF(SECOND, call.start_time, call.end_time),0)) > (SELECT AVG(DATEDIFF(SECOND, call.start_time, call.end_time)) FROM call)
ORDER BY calls DESC, country.id ASC;
`

	simplebad := "SELECT * FROM GROUP BY age"

	for _,  test := range []string{simple, medium, complex, simplebad} {
		ast, err := parse.Parse(sql.NewEmptyContext(), test)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(pretty.Sprint(ast))
	}
}
