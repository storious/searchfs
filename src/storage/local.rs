use std::fs;
use std::io;
use std::path::PathBuf;

use crate::storage::Storage;

pub struct LocalStorage {
    root: PathBuf,
}

impl LocalStorage {
    pub fn new<P: Into<PathBuf>>(root: P) -> Self {
        Self { root: root.into() }
    }
    fn full_path(&self, path: &str) -> PathBuf {
        self.root.join(path)
    }
}

impl Storage for LocalStorage {
    fn create_dir_all(&self, path: &str) -> io::Result<()> {
        fs::create_dir_all(self.full_path(path))
    }

    fn write(&self, path: &str, data: &[u8]) -> io::Result<()> {
        let full = self.full_path(path);

        if let Some(parent) = full.parent() {
            fs::create_dir_all(parent)?;
        }

        fs::write(full, data)
    }

    fn read(&self, path: &str) -> io::Result<Vec<u8>> {
        fs::read(self.full_path(path))
    }

    fn remove_dir_all(&self, path: &str) -> io::Result<()> {
        fs::remove_dir_all(self.full_path(path))
    }

    fn exists(&self, path: &str) -> bool {
        self.full_path(path).exists()
    }
}
