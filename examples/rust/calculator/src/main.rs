mod parser;

use parser::{PigeonParser, Value};

fn calculate(input: &str) -> Result<i64, String> {
    match PigeonParser::parse(input).map_err(|error| error.to_string())? {
        Value::Int(value) => Ok(value),
        value => Err(format!("unexpected parser result: {value:?}")),
    }
}

fn main() {
    let input = std::env::args().skip(1).collect::<Vec<_>>().join(" ");
    if input.is_empty() {
        eprintln!("usage: pigeon-calculator-example 'EXPRESSION'");
        std::process::exit(2);
    }
    match calculate(&input) {
        Ok(value) => println!("{value}"),
        Err(error) => {
            eprintln!("{error}");
            std::process::exit(1);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::calculate;

    #[test]
    fn evaluates_precedence_and_parentheses() {
        assert_eq!(calculate("18 + 3 - 27 * (-18 / -3)").unwrap(), -141);
        assert_eq!(calculate("1 + (2 - 3) * 4 / 5").unwrap(), 1);
    }

    #[test]
    fn rejects_incomplete_expressions() {
        assert!(calculate("1 +").is_err());
    }
}
