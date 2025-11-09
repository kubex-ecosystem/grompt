#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Check for possible secrets in documentation files.
"""


import pathlib
import re
import sys


class DocSecretChecker:
    """
    A class to check documentation files for possible secrets.
    """

    bad = []
    SUSPECT = [
        r"sk-[A-Za-z0-9_-]{10,}",   # openai-like
        r"gsk_[A-Za-z0-9_-]{10,}",  # groq-like
        r"AIza[0-9A-Za-z_-]{10,}",  # google-ish
        r"ghp_[A-Za-z0-9]{20,}",    # github
    ]
    ALLOW_INLINE = "pragma: allowlist secret"

    def __init__(self):
        self.bad = []

    def is_doc(self, p: pathlib.Path) -> bool:
        """
        Check if the file is a documentation file by its extension.
        """
        return p.suffix in {".md", ".mdx", ".rst", ".txt", ".adoc", ".html"}

    def has_secret(self, fname: str, text: str) -> bool:
        """
        Check if a line contains a secret based on predefined patterns.
        """
        if self.ALLOW_INLINE in text:
            # Linha a linha: só marca as que não têm pragma acima
            lines = text.splitlines()
            for i, line in enumerate(lines):
                for pat in self.SUSPECT:
                    if re.search(pat, line):
                        # olha até 2 linhas acima por pragma
                        if not any(self.ALLOW_INLINE in lines[j] for j in range(max(0, i-2), i)):
                            self.bad.append(
                                f"{fname}:{i+1}: {line.strip()[:120]}")
        else:
            for pat in self.SUSPECT:
                if re.search(pat, text):
                    self.bad.append(
                        f"{fname}: contains suspect pattern and no '{self.ALLOW_INLINE}'")

    def check(self):
        """
        Main function to check files for secrets.
        """

        for fname in sys.argv[1:]:
            p = pathlib.Path(fname)
            if not p.exists() or p.is_dir() or not self.is_doc(p):
                continue
            text = p.read_text(encoding="utf-8", errors="ignore")
            if self.has_secret(fname, text):

                continue

        if self.bad:
            print("Docs contain possible secrets (no pragma):")
            print("\n".join(self.bad))
            sys.exit(1)


def main():
    """
    Main function to run the DocSecretChecker.
    """
    DocSecretChecker().check()


if __name__ == "__main__":
    main()
