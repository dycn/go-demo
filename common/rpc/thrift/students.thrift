namespace go  students

struct Student{
 1: i32 id,
 2: string name,
 3: Sex sex,
 4: i16 age,
}
enum Sex {
  DEFAULT = 1,
  MALE = 2,
  FEMALE = 3
}

typedef set<Student> All

service ClassMember {
    void Ping(1:i64 num),
    string GetClassName(),
    list<Student> List(),
    void Add(1: Student s),
    bool IsNameExist(1:string name),
    i16 Count(),
}