package configs

import "testing"

func TestConfig_Read(t *testing.T) {
	type fields struct {
		Server   Server
		Database Database
		Crawler  Crawler
		Logger   Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Server:   tt.fields.Server,
				Database: tt.fields.Database,
				Crawler:  tt.fields.Crawler,
				Logger:   tt.fields.Logger,
			}
			if err := c.Read(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
