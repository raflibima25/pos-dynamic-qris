/**
 * Format number as Indonesian Rupiah currency
 */
export function formatRupiah(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount)
}

/**
 * Format number as simple Rupiah (without currency symbol)
 */
export function formatRupiahSimple(amount: number): string {
  return new Intl.NumberFormat('id-ID').format(amount)
}