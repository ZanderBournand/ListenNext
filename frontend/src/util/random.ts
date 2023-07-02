export function getTimeSince(timestamp: Date): string {
    const now = new Date();
    const timeDiff = now.getTime() - timestamp.getTime();
  
    // Define the time intervals in milliseconds
    const minute = 60 * 1000;
    const hour = 60 * minute;
    const day = 24 * hour;
    const month = 30 * day;
    const year = 365 * day;
  
    // Calculate the difference in each time interval
    const minutes = Math.floor(timeDiff / minute);
    const hours = Math.floor(timeDiff / hour);
    const days = Math.floor(timeDiff / day);
    const months = Math.floor(timeDiff / month);
    const years = Math.floor(timeDiff / year);
  
    // Return the appropriate "since" time string
    if (years > 0) {
      return `${years} year${years > 1 ? 's' : ''} ago`;
    } else if (months > 0) {
      return `${months} month${months > 1 ? 's' : ''} ago`;
    } else if (days > 0) {
      return `${days} day${days > 1 ? 's' : ''} ago`;
    } else if (hours > 0) {
      return `${hours} hour${hours > 1 ? 's' : ''} ago`;
    } else if (minutes > 0) {
      return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
    } else {
      return 'Just now';
    }
  }
  