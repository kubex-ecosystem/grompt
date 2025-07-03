# **Contributing to Grompt**

Thank you for your interest in contributing to **[Grompt](https://github.com/rafa-mori/grompt)**! We are excited to have you as part of our community. This guide will help you get started and contribute effectively to the project.

---

## **How to Contribute**

There are several ways to contribute to [Grompt](https://github.com/rafa-mori/grompt):

1. **Report Issues**
   - Found bugs or issues in the code? Open an issue detailing the problem.
   - Include as much information as possible: steps to reproduce the issue, logs, Go version used, etc.

2. **Suggest Improvements**
   - Have an idea to improve the project? Share your suggestion by opening an issue with the `enhancement` tag.

3. **Submit Pull Requests**
   - Want to fix a bug or implement something new? Submit a pull request with your changes.

4. **Test and Review Code**
   - Help review pull requests from other contributors.
   - Run existing tests and validate if the proposed changes keep the system functional.

---

## **Getting Started**

### 1. **Clone the Repository**

```bash
git clone https://github.com/rafa-mori/grompt.git
cd grompt
```

### 2. **Set Up the Environment**

Make sure you have Go installed:
You can install Go using one of the following methods:

- [A super easy way to install Go](https://github.com/rafa-mori/go-installer)

  ```shell
    curl -sSfL 'https://raw.githubusercontent.com/faelmori/go-installer/refs/heads/main/go.sh' | bash
  ```

- [Convencional way to download Go](https://go.dev/dl/)

### 3. **Install Dependencies**

```bash
# Download the necessary packages
go mod download
```

### 4. **Run Tests**+

Before making changes, run the existing tests:

```bash
go test ./...
```

---

## **Creating a Pull Request**

### **1. Fork the Repository**

Create a fork of the project to your own GitHub.

### **2. Create a New Branch**

```bash
git checkout -b your-feature
```

### **3. Make Changes**

```bash
# Edit the code in your favorite editor
```

Make sure to follow the project's code conventions and best practices.

### **4. Add Tests (if applicable)**

Include test cases to validate the added functionality.

### **5. Run Tests**

Ensure all changes and tests are working:

```bash
go test ./...
```

### **6. Commit and Push**

```bash
git add .
git commit -m "Brief description of the change"
git push origin your-feature
```

### **7. Open the Pull Request**

Go to the original repository on GitHub and open a pull request explaining your changes.

---

## **Code Standards**

### **Code Style**

This project follows Go's code conventions. Some recommendations:

- Use `gofmt` to format the code:

```bash
gofmt -w .
```

- Name variables and functions clearly and descriptively.
- Break down long functions into smaller parts whenever possible.

### **Commits**

Commits should be clear and descriptive. Examples:

- `fix: fix bug in notification logic`
- `feat: add support for Slack notifier`

---

## **Best Practices**

1. **Be Respectful and Welcoming**  
   This is an open-source project for everyone. Respect other contributors and collaborate constructively.

2. **Document Your Changes**  
   Update the `README.md` or documentation, if necessary, to include your changes.

3. **Add Tests When Possible**  
   Ensure any new functionality is accompanied by tests.

4. **Be Clear in Issue Reports**  
   When opening an issue, be detailed and provide as much context as possible.

---

## **Where to Get Help**

If you need assistance, feel free to:

- Open an issue with the `question` tag.
- Contact me via the email or LinkedIn listed in the `README.md`.

---

## **Our Commitment**

We commit to reviewing pull requests and issues as quickly as possible. We value your contribution and appreciate the time dedicated to the project!
