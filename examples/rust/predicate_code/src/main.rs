mod parser;

use parser::{PTParser, Value};

fn parse_operator(input: &str) -> Result<String, String> {
    match PTParser::parse(input).map_err(|error| error.to_string())? {
        Value::String(value) => Ok(value),
        value => Err(format!("unexpected parser result: {value:?}")),
    }
}

fn main() {
    let input = std::env::args().nth(1).unwrap_or_else(|| "%".to_string());
    match parse_operator(&input) {
        Ok(value) => println!("{value}"),
        Err(error) => {
            eprintln!("{error}");
            std::process::exit(1);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::parse_operator;

    #[test]
    fn accepts_configured_operators() {
        assert_eq!(parse_operator("+").unwrap(), "+");
        assert_eq!(parse_operator("%").unwrap(), "%");
    }

    #[test]
    fn rejects_other_symbols() {
        assert!(parse_operator("-").is_err());
    }
}
