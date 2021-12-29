export const validatePassword = (password: string): boolean => {
  // ^                         Start anchor
  // (?=.*[A-Z])               Ensure string has one uppercase letter.
  // (?=.*[!@#$&*])            Ensure string has one special case letter.
  // (?=.*[0-9])               Ensure string has one digit.
  // (?=.*[a-z])               Ensure string has one lowercase letter.
  // .{8}                      Ensure string is of length 8.
  // $                         End anchor.
  // const strictSecureRegexp = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%\^&\*])(?=.{8,})/
  const securePasswordRegexp = /^(((?=.*[a-z])(?=.*[A-Z]))|((?=.*[a-z])(?=.*[0-9]))|((?=.*[A-Z])(?=.*[0-9])))(?=.{6,})/
  return password.match(securePasswordRegexp) !== null
}