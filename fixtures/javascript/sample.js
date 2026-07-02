// Simple function
function simple() {
    return 42;
}

// Complex function
function complexFunc(x, y) {
    if (x > 0) {
        for (let i = 0; i < x; i++) {
            if (i % 2 === 0) {
                y += i;
            } else if (i % 3 === 0) {
                y -= i;
            } else {
                y *= 2;
            }
        }
    } else if (x < 0) {
        while (y > 0) {
            y--;
        }
    } else {
        try {
            y = x / y;
        } catch (e) {
            y = 0;
        }
    }
    return y;
}

// Arrow function with ternary
const classify = (x) => x > 100 ? "large" : x > 10 ? "medium" : "small";

// Class
class Calculator {
    add(a, b) {
        return a + b;
    }

    compute(operation, a, b) {
        switch (operation) {
            case "add":
                return this.add(a, b);
            case "sub":
                return a - b;
            case "mul":
                return a * b;
            case "div":
                if (b !== 0) {
                    return a / b;
                }
                throw new Error("Division by zero");
            default:
                throw new Error(`Unknown operation: ${operation}`);
        }
    }
}

module.exports = { simple, complexFunc, classify, Calculator };
