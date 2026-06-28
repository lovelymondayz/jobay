export const convertToWIBTimeString = (dateStr?: string): string => {
  if (!dateStr) return "";

  const date = new Date(dateStr);
  const currentOffset = date.getTimezoneOffset(); // in minutes (negative means ahead of UTC)
  const wibOffset = -420; // UTC+7

  // Jika offset saat ini bukan WIB, sesuaikan
  if (currentOffset !== wibOffset) {
    // Konversi ke WIB dengan cara menyesuaikan jamnya
    date.setMinutes(date.getMinutes() + (currentOffset - wibOffset));
  }

  // Format ke "YYYY-MM-DD HH:mm:ss"
  const pad = (n: number) => n.toString().padStart(2, "0");

  const year = date.getFullYear();
  const month = pad(date.getMonth() + 1);
  const day = pad(date.getDate());
  const hour = pad(date.getHours());
  const minute = pad(date.getMinutes());
  const second = pad(date.getSeconds());

  return `${year}-${month}-${day} ${hour}:${minute}:${second}`;
};
