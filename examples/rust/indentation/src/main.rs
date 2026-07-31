mod parser;

use parser::{execute, parse_program};

fn run(input: &str) -> Result<i64, String> {
    let program = parse_program(input).map_err(|error| error.to_string())?;
    execute(&program)
}

fn main() {
    let path = std::env::args().nth(1).unwrap_or_else(|| {
        eprintln!("usage: pegtool-indentation-example FILE");
        std::process::exit(2);
    });
    let input = std::fs::read_to_string(path).unwrap_or_else(|error| {
        eprintln!("{error}");
        std::process::exit(1);
    });
    match run(&input) {
        Ok(value) => println!("{value}"),
        Err(error) => {
            eprintln!("{error}");
            std::process::exit(1);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::run;

    #[test]
    fn executes_significant_whitespace() {
        let input = "x = 40\nif x:\n    x = x + 2\nreturn x\n";
        assert_eq!(run(input).unwrap(), 42);
    }

    #[test]
    fn skips_false_if_blocks() {
        let input = "x = 0\nif x:\n    x = 5\nreturn x\n";
        assert_eq!(run(input).unwrap(), 0);
    }

    #[test]
    fn rejects_bad_indentation() {
        let input = "x = 1\nif x:\n  x = 2\nreturn x\n";
        assert!(run(input).is_err());
    }
}
