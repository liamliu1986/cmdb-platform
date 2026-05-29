package config

import (
    "github.com/spf13/viper"
    "log"
)

type Config struct {
    ServerPort     string `mapstructure:"SERVER_PORT"`
    DBHost         string `mapstructure:"DB_HOST"`
    DBPort         string `mapstructure:"DB_PORT"`
    DBUser         string `mapstructure:"DB_USER"`
    DBPassword     string `mapstructure:"DB_PASSWORD"`
    DBName         string `mapstructure:"DB_NAME"`
    RedisHost      string `mapstructure:"REDIS_HOST"`
    RedisPort      string `mapstructure:"REDIS_PORT"`
    RedisPassword  string `mapstructure:"REDIS_PASSWORD"`
    RedisDB        int    `mapstructure:"REDIS_DB"`
    JWTSecret      string `mapstructure:"JWT_SECRET"`
    JWTExpireHours int    `mapstructure:"JWT_EXPIRE_HOURS"`
}

func Load() *Config {
    viper.AutomaticEnv()
    viper.SetDefault("SERVER_PORT", "8080")
    viper.SetDefault("DB_HOST", "localhost")
    viper.SetDefault("DB_PORT", "5432")
    viper.SetDefault("DB_USER", "cmdb")
    viper.SetDefault("DB_NAME", "cmdb")
    viper.SetDefault("REDIS_HOST", "localhost")
    viper.SetDefault("REDIS_PORT", "6379")
    viper.SetDefault("REDIS_DB", 0)
    viper.SetDefault("JWT_EXPIRE_HOURS", 24)

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        log.Fatal("Failed to load config:", err)
    }
    return &cfg
}
