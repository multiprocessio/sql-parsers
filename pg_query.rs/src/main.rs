use pg_query;

fn main() {
    let simple = "SELECT * FROM x";

    let medium = "SELECT COUNT(1) AS count, name section FROM t GROUP BY name LIMIT 10";

    let complex = "
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
";

    let simplebad = "SELECT * FROM GROUP BY age";

    let tests = vec![simple, medium, complex, simplebad];

    for test in tests {
        match pg_query::parse(test) {
            Ok(ast) => println!("AST: {:#?}", ast),
            Err(err) => println!("Error: {:#?}", err),
        }
    }
}
