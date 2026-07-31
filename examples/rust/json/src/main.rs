mod parser;

use parser::{PTParser, Value};

fn parse_json(input: &str) -> Result<Value, String> {
    PTParser::parse(input).map_err(|error| error.to_string())
}

fn main() {
    let input = std::env::args()
        .nth(1)
        .unwrap_or_else(|| r#"{"pigeon": true, "targets": ["go", "rust"]}"#.to_string());
    match parse_json(&input) {
        Ok(value) => println!("{value:?}"),
        Err(error) => {
            eprintln!("{error}");
            std::process::exit(1);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::{parse_json, Value};

    #[test]
    fn parses_objects_arrays_and_null() {
        let value = parse_json(r#"{"items": [1, null, true], "name": "pigeon"}"#).unwrap();
        let Value::Object(object) = value else {
            panic!("expected object");
        };
        assert!(matches!(object.get("name"), Some(Value::String(value)) if value == "pigeon"));
        let Some(Value::Array(items)) = object.get("items") else {
            panic!("expected items array");
        };
        assert_eq!(items.len(), 3);
        assert!(matches!(items[0], Value::Float(value) if value == 1.0));
        assert!(matches!(items[1], Value::Null));
        assert!(matches!(items[2], Value::Bool(true)));
    }

    #[test]
    fn decodes_escapes_and_surrogate_pairs() {
        let value = parse_json(r#""line\n\uD83D\uDE00\/""#).unwrap();
        assert!(matches!(value, Value::String(value) if value == "line\n😀/"));
    }

    #[test]
    fn rejects_leading_zeroes() {
        assert!(parse_json("00").is_err());
    }
}
