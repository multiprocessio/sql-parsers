package main;

import com.google.gson.*;
import com.facebook.presto.sql.parser.*;

public class Main {
    public static void main() {
	var simple = "SELECT * FROM x";
	
	var medium = "SELECT COUNT(1) AS count, name section FROM t GROUP BY name LIMIT 10";
	
	var complex = """
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
	""";

	var simplebad = "SELECT * FROM GROUP BY age";

	var parser = new SqlParser();
	String[] tests = {simple, medium, complex, simplebad};
	for (var test : tests) {
	    var query = (Query)parser.createStatement(sql);
	    System.out.println(gson.toJson(query).toString());
	}
    }
}
