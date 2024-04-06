def divide_names_by_letter(filename):

  with open(filename, "r") as f:
    names = [line.strip() for line in f]

  # Create a dictionary to store names by first letter
  name_dict = {}
  for name in names:
    first_letter = name[0].upper()  # Use uppercase for case-insensitive sorting
    if first_letter not in name_dict:
      name_dict[first_letter] = []
    name_dict[first_letter].append(name)

  # Write names to separate files
  for letter, names in name_dict.items():
    with open(f"male_{letter}.txt", "w") as f:
      f.write("\n".join(names))

  print("Names divided into separate files by first letter!")

# Example usage
filename = "mnames.txt"  # Replace with your actual file name
divide_names_by_letter(filename)

