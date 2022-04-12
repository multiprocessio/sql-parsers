open Sqlgg

let simple = "SELECT * FROM x"

let medium = "SELECT COUNT(1) AS count, name section FROM t GROUP BY name LIMIT 10"

let complex = {|
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
ORDER BY talks DESC, country.id ASC
|}

let simplebad = "SELECT * FROM GROUP BY age"

let tests = [simple, true; medium, true; complex, true; simplebad, false]


let () =
  Format.printf "@[<v>";
  let failed =
    tests |>
    List.fold_left (fun failed (sql, expected) ->
      let actual =
        try
          begin match Parser.T.parse_buf_exn (Lexing.from_string sql) with
          | Select sf -> Format.printf "%a@;" Sql.pp_select_full sf; true
          | _ -> failwith "expected a SELECT"
          end
        (* adapted from src/main.ml *)
        with Parser_utils.Error (exn,(line,cnum,tok,tail)) ->
          let sql = String.trim sql in
          let extra =
            match exn with
            | Sql.Schema.Error (_,msg) -> msg
            | exn -> Printexc.to_string exn
          in
          if cnum = String.length sql && tok = "" then
            Format.printf "%s@;Error: %s@;" sql extra
          else (
            let t = String.(sub tail 0 (min (length tail) 32)) in
            Format.printf "%s@;Position %u:%u Tokens: %s%s@;Error: %s@;" sql line cnum tok t extra
          );
          false
      in
      failed || actual <> expected
    ) false
  in
  Format.printf "@]@.";
  if failed then exit 1
