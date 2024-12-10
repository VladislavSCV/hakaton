package utils

//func saveImage(userID, imagePath string) error {
//	fileName := filepath.Base(imagePath)
//	destPath := fmt.Sprintf("uploads/%s/%s", userID, fileName)
//
//	err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
//	if err != nil {
//		return err
//	}
//
//	srcFile, err := os.Open(imagePath)
//	if err != nil {
//		return err
//	}
//	defer srcFile.Close()
//
//	destFile, err := os.Create(destPath)
//	if err != nil {
//		return err
//	}
//	defer destFile.Close()
//
//	_, err = io.Copy(destFile, srcFile)
//	if err != nil {
//		return err
//	}
//
//	return repository.SaveImage(userID, destPath)
//}
