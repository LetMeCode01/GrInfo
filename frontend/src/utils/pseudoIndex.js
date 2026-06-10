export function pseudoIndexFromId(id, mod = 17) {
  let sum = 0;
  for (let i = 0; i < String(id).length; i++) sum += String(id).charCodeAt(i);
  return sum % mod;
}
