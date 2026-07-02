# Simple function
def simple():
    return 42

# Function with decision points
def complex_func(x, y):
    if x > 0:
        for i in range(x):
            if i % 2 == 0:
                y += i
            elif i % 3 == 0:
                y -= i
            else:
                y *= 2
    elif x < 0:
        while y > 0:
            y -= 1
    else:
        try:
            y = x / y
        except ZeroDivisionError:
            y = 0
    return y

# Function with deep nesting
def deeply_nested(data):
    result = []
    for item in data:
        if item is not None:
            if isinstance(item, list):
                for sub in item:
                    if sub > 0:
                        result.append(sub)
            elif isinstance(item, dict):
                for key, value in item.items():
                    if value:
                        result.append(key)
    return result

# Class with methods
class Calculator:
    def add(self, a, b):
        return a + b

    def compute(self, operation, a, b):
        if operation == "add":
            return self.add(a, b)
        elif operation == "sub":
            return a - b
        elif operation == "mul":
            return a * b
        elif operation == "div":
            if b != 0:
                return a / b
            else:
                raise ValueError("Division by zero")
        else:
            raise ValueError(f"Unknown operation: {operation}")
