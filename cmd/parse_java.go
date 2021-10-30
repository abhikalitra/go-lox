package main

import (
	"lisp/java"
	"log"
)

func main() {
	s := `
package io.synlabs.projektor.entity;

import io.synlabs.projektor.enums.ClientType;
import lombok.Getter;
import lombok.Setter;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.EnumType;
import javax.persistence.Enumerated;

@Entity
@Getter
@Setter
public class Client extends BaseEntity {

    @Column(nullable = false, length = 64)
    private String name;

    @Column(nullable = false, length = 64)
    private String location;

    @Enumerated(EnumType.STRING)
    private ClientType clientType;
}
	`
	log.Println("hello")
	js := java.NewScanner()
	js.Eval(s)
	//for _, tok := range js.Tokens {
	//	log.Printf("%+v", tok)
	//}

	jp := java.NewParser(js.Tokens)
	statements := jp.Parse()

	for _, stmt := range statements {
		stmt.PrintStmt()
	}

}
