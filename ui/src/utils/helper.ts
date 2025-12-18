export const isRTL = (text: string): boolean => {
  if (!text) return false;
  const rtl = text.match(/[\u0600-\u08FF\uFB50-\uFDFF\uFE70-\uFEFF]/g) || []
  const ltr = text.match(/[A-Za-z]/g) || []
  return rtl.length > 0 && ltr.length === 0 || rtl.length >= ltr.length * 2
}