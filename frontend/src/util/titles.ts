export function ReduceName(name: any, max: number) {
    if (name.length <= max) {
      return name; // Return the name as is if it is 20 characters or less
    } else {
      return name.substring(0, max) + '...'; // Add ellipsis if the name exceeds 20 characters
    }
  }