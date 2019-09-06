package pkg

// ConnectDatabase -
func ConnectDatabase(driver DatabaseDriver) error {
	driver.Connect();
	return nil
}
