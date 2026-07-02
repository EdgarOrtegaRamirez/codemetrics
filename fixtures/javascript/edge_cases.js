// Edge cases for JavaScript testing

// Simple arrow function
const add = (a, b) => a + b;

// Arrow function with block body
const multiply = (a, b) => {
  return a * b;
};

// Function declaration
function subtract(a, b) {
  return a - b;
}

// Function expression
const divide = function(a, b) {
  if (b === 0) {
    throw new Error('Division by zero');
  }
  return a / b;
};

// Async function
async function fetchData(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Fetch failed:', error);
    return null;
  }
}

// Generator function
function* range(start, end) {
  for (let i = start; i < end; i++) {
    yield i;
  }
}

// Class with methods
class Calculator {
  constructor() {
    this.result = 0;
  }

  add(value) {
    this.result += value;
    return this;
  }

  subtract(value) {
    this.result -= value;
    return this;
  }

  getResult() {
    return this.result;
  }
}

// IIFE
const singleton = (function() {
  let instance;
  function createInstance() {
    return { name: 'singleton' };
  }
  return {
    getInstance: function() {
      if (!instance) {
        instance = createInstance();
      }
      return instance;
    }
  };
})();

// Higher-order function
const filter = (arr, predicate) => {
  const result = [];
  for (const item of arr) {
    if (predicate(item)) {
      result.push(item);
    }
  }
  return result;
};

// Closure
function counter() {
  let count = 0;
  return {
    increment: () => ++count,
    decrement: () => --count,
    getCount: () => count
  };
}

// Complex nested logic
function processOrder(order) {
  if (!order) {
    return { error: 'No order' };
  }

  let total = 0;
  const discounts = [];

  for (const item of order.items) {
    if (item.quantity > 10) {
      discounts.push({ item: item.name, discount: 0.1 });
      total += item.price * item.quantity * 0.9;
    } else if (item.quantity > 5) {
      discounts.push({ item: item.name, discount: 0.05 });
      total += item.price * item.quantity * 0.95;
    } else {
      total += item.price * item.quantity;
    }
  }

  if (total > 1000) {
    total *= 0.95; // 5% bulk discount
  }

  return { total, discounts };
}
