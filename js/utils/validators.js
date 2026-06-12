/**
 * Input validation functions
 */

/**
 * Validate email
 */
function validateEmail(email) {
  if (!email) return { valid: false, error: 'Email is required' };
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!re.test(email)) return { valid: false, error: 'Invalid email format' };
  return { valid: true };
}

/**
 * Validate password
 */
function validatePassword(password) {
  if (!password) return { valid: false, error: 'Password is required' };
  if (password.length < 8) return { valid: false, error: 'Mật khẩu phải có ít nhất 8 ký tự' };
  return { valid: true };
}

/**
 * Validate phone number
 */
function validatePhone(phone) {
  if (!phone) return { valid: false, error: 'Phone is required' };
  const re = /^(0|84)[0-9]{9,10}$/;
  if (!re.test(phone.replace(/\s/g, ''))) {
    return { valid: false, error: 'Invalid phone number format' };
  }
  return { valid: true };
}

/**
 * Validate required field
 */
function validateRequired(value, fieldName) {
  if (!value || (typeof value === 'string' && value.trim() === '')) {
    return { valid: false, error: `${fieldName} is required` };
  }
  return { valid: true };
}

/**
 * Validate number
 */
function validateNumber(value, fieldName, min = 0, max = null) {
  if (!value && value !== 0) return { valid: false, error: `${fieldName} is required` };
  const num = parseFloat(value);
  if (isNaN(num)) return { valid: false, error: `${fieldName} must be a number` };
  if (num < min) return { valid: false, error: `${fieldName} must be at least ${min}` };
  if (max !== null && num > max) {
    return { valid: false, error: `${fieldName} must not exceed ${max}` };
  }
  return { valid: true };
}

/**
 * Validate date
 */
function validateDate(date, fieldName) {
  if (!date) return { valid: false, error: `${fieldName} is required` };
  const d = new Date(date);
  if (isNaN(d.getTime())) return { valid: false, error: `${fieldName} is invalid` };
  return { valid: true };
}

/**
 * Validate form object
 */
function validateForm(formData, rules) {
  const errors = {};
  for (const field in rules) {
    const value = formData[field];
    const fieldRules = rules[field];

    for (const rule of fieldRules) {
      if (rule.type === 'required') {
        const result = validateRequired(value, rule.label || field);
        if (!result.valid) {
          errors[field] = result.error;
          break;
        }
      } else if (rule.type === 'email') {
        const result = validateEmail(value);
        if (!result.valid) {
          errors[field] = result.error;
          break;
        }
      } else if (rule.type === 'password') {
        const result = validatePassword(value);
        if (!result.valid) {
          errors[field] = result.error;
          break;
        }
      } else if (rule.type === 'phone') {
        const result = validatePhone(value);
        if (!result.valid) {
          errors[field] = result.error;
          break;
        }
      } else if (rule.type === 'number') {
        const result = validateNumber(value, rule.label || field, rule.min, rule.max);
        if (!result.valid) {
          errors[field] = result.error;
          break;
        }
      } else if (rule.type === 'date') {
        const result = validateDate(value, rule.label || field);
        if (!result.valid) {
          errors[field] = result.error;
          break;
        }
      } else if (rule.type === 'custom' && rule.validate) {
        const result = rule.validate(value);
        if (!result.valid) {
          errors[field] = result.error;
          break;
        }
      }
    }
  }
  return { valid: Object.keys(errors).length === 0, errors };
}

export {
  validateEmail,
  validatePassword,
  validatePhone,
  validateRequired,
  validateNumber,
  validateDate,
  validateForm,
};
