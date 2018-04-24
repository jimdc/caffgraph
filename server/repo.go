package main

import "fmt"

var currentId int
var doses []Dose

// Give us some seed data
func init() {
    RepoCreateDose(Dose{Name:"Pill", Dosage: 200,
                   Time: "2018-04-13T12:00:00Z", Remnants: nil})

    RepoCreateDose(Dose{Name: "Dunkin 13.7FlOz", Dosage: 198,
                   Time: "2018-14-14T15:15:00Z", Remnants: nil})
}

func RepoFindDose(id int) Dose {
    for _, d := range doses {
        if d.Id == id {
            return d
        }
    }
    // return empty Dose if not found
    return Dose{}
}

func RepoCreateDose(d Dose) Dose {
    currentId += 1
    d.Id = currentId
    doses = append(doses, d)
    return d
}

func RepoDestroyDose(id int) error {
    for i, d := range doses {
        if d.Id == id {
            doses = append(doses[:i], doses[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Dose with id of %d to delete", id)
}

