// Simple function
fn simple() -> i32 {
    42
}

// Complex function
fn complex_func(x: i32, mut y: i32) -> i32 {
    if x > 0 {
        for i in 0..x {
            if i % 2 == 0 {
                y += i;
            } else if i % 3 == 0 {
                y -= i;
            } else {
                y *= 2;
            }
        }
    } else if x < 0 {
        while y > 0 {
            y -= 1;
        }
    } else {
        if y != 0 {
            y = x / y;
        } else {
            y = 0;
        }
    }
    y
}

// Function with match
fn classify(x: i32) -> &'static str {
    match x {
        x if x > 100 => "large",
        x if x > 10 => "medium",
        x if x > 0 => "small",
        _ => "negative",
    }
}

// Struct with methods
struct Calculator;

impl Calculator {
    fn compute(&self, op: &str, a: i32, b: i32) -> Result<i32, String> {
        match op {
            "add" => Ok(a + b),
            "sub" => Ok(a - b),
            "mul" => Ok(a * b),
            "div" => {
                if b != 0 {
                    Ok(a / b)
                } else {
                    Err("Division by zero".to_string())
                }
            }
            _ => Err(format!("Unknown operation: {}", op)),
        }
    }
}

fn main() {
    println!("{}", simple());
}
