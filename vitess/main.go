package main

import (
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"
	"github.com/kr/pretty"
)

func main() {
	simple := "SELECT * FROM x"
	
	medium := "SELECT COUNT(1) AS count, name section FROM t GROUP BY name LIMIT 10"
	
	complex := `
SELECT 
	country.country_name_eng,
	SUM(CASE WHEN talk.id IS NOT NULL THEN 1 ELSE 0 END) AS talks,
	AVG(ISNULL(DATEDIFF(SECOND, talk.start_time, talk.end_time),0)) AS avg_difference
FROM country 
LEFT JOIN city ON city.country_id = country.id
LEFT JOIN customer ON city.id = customer.city_id
LEFT JOIN talk ON talk.customer_id = customer.id
GROUP BY 
	country.id,
	country.country_name_eng
HAVING AVG(ISNULL(DATEDIFF(SECOND, talk.start_time, talk.end_time),0)) > (SELECT AVG(DATEDIFF(SECOND, talk.start_time, talk.end_time)) FROM talk)
ORDER BY talks DESC, country.id ASC;
`

	simplebad := "SELECT * FROM GROUP BY age"

	for _,  test := range []string{simple, medium, complex, simplebad} {
		ast, _, err := sqlparser.Parse2(test)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(pretty.Sprint(ast))
	}
}
