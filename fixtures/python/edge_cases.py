# Edge cases for testing
def no_docstring():
    pass

def single_line():
    return 42

def with_defaults(a, b=10, c="hello"):
    return a + b

def nested_default(data=None):
    if data is None:
        data = []
    return data

# Test try/except/finally
def error_handling():
    try:
        result = 1 / 0
    except ZeroDivisionError:
        result = 0
    except Exception as e:
        result = -1
    finally:
        pass
    return result

# Test list comprehension
def list_comp():
    return [x for x in range(10) if x % 2 == 0]

# Test generator
def generator():
    for i in range(10):
        yield i

# Test lambda
add = lambda x, y: x + y

# Test nested functions
def outer():
    def inner():
        return 42
    return inner()

# Test class with inheritance
class Base:
    def method(self):
        pass

class Derived(Base):
    def method(self):
        super().method()
        return True
