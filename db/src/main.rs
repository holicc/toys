use std::io;
use std::io::Write;
use std::process::exit;

enum FLAG {
    EXIT,
    HELP,
    QUERY(String),
}

struct Repl {
    i: io::Stdin,
    o: io::Stdout,
}

fn main() {
    //TODO engine server

    //
    let mut repl = Repl::new();
    loop {
        repl.print("db:> ");
        match repl.read_line() {
            FLAG::EXIT => {
                repl.print("Bye!");
                exit(0);
            }
            FLAG::HELP => {
                repl.print("help message here. \n");
            }
            FLAG::QUERY(query) => {
                if query.starts_with("select") {
                    repl.print("do select operation \n");
                } else if query.starts_with("insert") {
                    repl.print("do insert operation \n");
                } else {
                    //TODO BUG display bug
                    //eg: ' recognized keyword at start of 'sasd
                    repl.print(format!("Unrecognized keyword at start of '{}' \n", query.as_str()).as_str())
                }
            }
        }
    }
}


impl Repl {
    fn new() -> Repl {
        return Repl {
            i: io::stdin(),
            o: io::stdout(),
        };
    }

    fn print(&mut self, s: &str) {
        self.o.write(s.as_bytes()).unwrap();
        self.o.flush().unwrap();
    }

    fn read_line(&self) -> FLAG {
        let mut r = String::new();
        self.i.read_line(&mut r).unwrap();
        match r.trim_end_matches(|x| x == '\n') {
            "\\q" => FLAG::EXIT,
            "\\h" => FLAG::HELP,
            query => FLAG::QUERY(String::from(query).to_lowercase())
        }
    }
}

