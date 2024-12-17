// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: proto/education.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	EducationService_GetCourses_FullMethodName    = "/education.EducationService/GetCourses"
	EducationService_GetCourseByID_FullMethodName = "/education.EducationService/GetCourseByID"
	EducationService_CreateCourse_FullMethodName  = "/education.EducationService/CreateCourse"
	EducationService_UpdateCourse_FullMethodName  = "/education.EducationService/UpdateCourse"
	EducationService_DeleteCourse_FullMethodName  = "/education.EducationService/DeleteCourse"
)

// EducationServiceClient is the client API for EducationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EducationServiceClient interface {
	GetCourses(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*CourseList, error)
	GetCourseByID(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*Course, error)
	CreateCourse(ctx context.Context, in *NewCourseRequest, opts ...grpc.CallOption) (*Course, error)
	UpdateCourse(ctx context.Context, in *UpdateCourseRequest, opts ...grpc.CallOption) (*Course, error)
	DeleteCourse(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*Empty, error)
}

type educationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEducationServiceClient(cc grpc.ClientConnInterface) EducationServiceClient {
	return &educationServiceClient{cc}
}

func (c *educationServiceClient) GetCourses(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*CourseList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseList)
	err := c.cc.Invoke(ctx, EducationService_GetCourses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *educationServiceClient) GetCourseByID(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*Course, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Course)
	err := c.cc.Invoke(ctx, EducationService_GetCourseByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *educationServiceClient) CreateCourse(ctx context.Context, in *NewCourseRequest, opts ...grpc.CallOption) (*Course, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Course)
	err := c.cc.Invoke(ctx, EducationService_CreateCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *educationServiceClient) UpdateCourse(ctx context.Context, in *UpdateCourseRequest, opts ...grpc.CallOption) (*Course, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Course)
	err := c.cc.Invoke(ctx, EducationService_UpdateCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *educationServiceClient) DeleteCourse(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, EducationService_DeleteCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EducationServiceServer is the server API for EducationService service.
// All implementations must embed UnimplementedEducationServiceServer
// for forward compatibility.
type EducationServiceServer interface {
	GetCourses(context.Context, *Empty) (*CourseList, error)
	GetCourseByID(context.Context, *CourseIDRequest) (*Course, error)
	CreateCourse(context.Context, *NewCourseRequest) (*Course, error)
	UpdateCourse(context.Context, *UpdateCourseRequest) (*Course, error)
	DeleteCourse(context.Context, *CourseIDRequest) (*Empty, error)
	mustEmbedUnimplementedEducationServiceServer()
}

// UnimplementedEducationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEducationServiceServer struct{}

func (UnimplementedEducationServiceServer) GetCourses(context.Context, *Empty) (*CourseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCourses not implemented")
}
func (UnimplementedEducationServiceServer) GetCourseByID(context.Context, *CourseIDRequest) (*Course, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCourseByID not implemented")
}
func (UnimplementedEducationServiceServer) CreateCourse(context.Context, *NewCourseRequest) (*Course, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCourse not implemented")
}
func (UnimplementedEducationServiceServer) UpdateCourse(context.Context, *UpdateCourseRequest) (*Course, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCourse not implemented")
}
func (UnimplementedEducationServiceServer) DeleteCourse(context.Context, *CourseIDRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCourse not implemented")
}
func (UnimplementedEducationServiceServer) mustEmbedUnimplementedEducationServiceServer() {}
func (UnimplementedEducationServiceServer) testEmbeddedByValue()                          {}

// UnsafeEducationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EducationServiceServer will
// result in compilation errors.
type UnsafeEducationServiceServer interface {
	mustEmbedUnimplementedEducationServiceServer()
}

func RegisterEducationServiceServer(s grpc.ServiceRegistrar, srv EducationServiceServer) {
	// If the following call pancis, it indicates UnimplementedEducationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EducationService_ServiceDesc, srv)
}

func _EducationService_GetCourses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EducationServiceServer).GetCourses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EducationService_GetCourses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EducationServiceServer).GetCourses(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _EducationService_GetCourseByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EducationServiceServer).GetCourseByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EducationService_GetCourseByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EducationServiceServer).GetCourseByID(ctx, req.(*CourseIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EducationService_CreateCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewCourseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EducationServiceServer).CreateCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EducationService_CreateCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EducationServiceServer).CreateCourse(ctx, req.(*NewCourseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EducationService_UpdateCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCourseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EducationServiceServer).UpdateCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EducationService_UpdateCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EducationServiceServer).UpdateCourse(ctx, req.(*UpdateCourseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EducationService_DeleteCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EducationServiceServer).DeleteCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EducationService_DeleteCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EducationServiceServer).DeleteCourse(ctx, req.(*CourseIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EducationService_ServiceDesc is the grpc.ServiceDesc for EducationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EducationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "education.EducationService",
	HandlerType: (*EducationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCourses",
			Handler:    _EducationService_GetCourses_Handler,
		},
		{
			MethodName: "GetCourseByID",
			Handler:    _EducationService_GetCourseByID_Handler,
		},
		{
			MethodName: "CreateCourse",
			Handler:    _EducationService_CreateCourse_Handler,
		},
		{
			MethodName: "UpdateCourse",
			Handler:    _EducationService_UpdateCourse_Handler,
		},
		{
			MethodName: "DeleteCourse",
			Handler:    _EducationService_DeleteCourse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/education.proto",
}

const (
	StudentService_RegisterStudent_FullMethodName      = "/education.StudentService/RegisterStudent"
	StudentService_LoginStudent_FullMethodName         = "/education.StudentService/LoginStudent"
	StudentService_GetStudentProfile_FullMethodName    = "/education.StudentService/GetStudentProfile"
	StudentService_UpdateStudentProfile_FullMethodName = "/education.StudentService/UpdateStudentProfile"
)

// StudentServiceClient is the client API for StudentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StudentServiceClient interface {
	RegisterStudent(ctx context.Context, in *RegisterStudentRequest, opts ...grpc.CallOption) (*Student, error)
	LoginStudent(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	GetStudentProfile(ctx context.Context, in *StudentIDRequest, opts ...grpc.CallOption) (*Student, error)
	UpdateStudentProfile(ctx context.Context, in *UpdateStudentRequest, opts ...grpc.CallOption) (*Student, error)
}

type studentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStudentServiceClient(cc grpc.ClientConnInterface) StudentServiceClient {
	return &studentServiceClient{cc}
}

func (c *studentServiceClient) RegisterStudent(ctx context.Context, in *RegisterStudentRequest, opts ...grpc.CallOption) (*Student, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Student)
	err := c.cc.Invoke(ctx, StudentService_RegisterStudent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentServiceClient) LoginStudent(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, StudentService_LoginStudent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentServiceClient) GetStudentProfile(ctx context.Context, in *StudentIDRequest, opts ...grpc.CallOption) (*Student, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Student)
	err := c.cc.Invoke(ctx, StudentService_GetStudentProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentServiceClient) UpdateStudentProfile(ctx context.Context, in *UpdateStudentRequest, opts ...grpc.CallOption) (*Student, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Student)
	err := c.cc.Invoke(ctx, StudentService_UpdateStudentProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StudentServiceServer is the server API for StudentService service.
// All implementations must embed UnimplementedStudentServiceServer
// for forward compatibility.
type StudentServiceServer interface {
	RegisterStudent(context.Context, *RegisterStudentRequest) (*Student, error)
	LoginStudent(context.Context, *LoginRequest) (*AuthResponse, error)
	GetStudentProfile(context.Context, *StudentIDRequest) (*Student, error)
	UpdateStudentProfile(context.Context, *UpdateStudentRequest) (*Student, error)
	mustEmbedUnimplementedStudentServiceServer()
}

// UnimplementedStudentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStudentServiceServer struct{}

func (UnimplementedStudentServiceServer) RegisterStudent(context.Context, *RegisterStudentRequest) (*Student, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterStudent not implemented")
}
func (UnimplementedStudentServiceServer) LoginStudent(context.Context, *LoginRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginStudent not implemented")
}
func (UnimplementedStudentServiceServer) GetStudentProfile(context.Context, *StudentIDRequest) (*Student, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentProfile not implemented")
}
func (UnimplementedStudentServiceServer) UpdateStudentProfile(context.Context, *UpdateStudentRequest) (*Student, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStudentProfile not implemented")
}
func (UnimplementedStudentServiceServer) mustEmbedUnimplementedStudentServiceServer() {}
func (UnimplementedStudentServiceServer) testEmbeddedByValue()                        {}

// UnsafeStudentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StudentServiceServer will
// result in compilation errors.
type UnsafeStudentServiceServer interface {
	mustEmbedUnimplementedStudentServiceServer()
}

func RegisterStudentServiceServer(s grpc.ServiceRegistrar, srv StudentServiceServer) {
	// If the following call pancis, it indicates UnimplementedStudentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StudentService_ServiceDesc, srv)
}

func _StudentService_RegisterStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterStudentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).RegisterStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_RegisterStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).RegisterStudent(ctx, req.(*RegisterStudentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentService_LoginStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).LoginStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_LoginStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).LoginStudent(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentService_GetStudentProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).GetStudentProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_GetStudentProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).GetStudentProfile(ctx, req.(*StudentIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentService_UpdateStudentProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStudentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceServer).UpdateStudentProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentService_UpdateStudentProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceServer).UpdateStudentProfile(ctx, req.(*UpdateStudentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StudentService_ServiceDesc is the grpc.ServiceDesc for StudentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StudentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "education.StudentService",
	HandlerType: (*StudentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterStudent",
			Handler:    _StudentService_RegisterStudent_Handler,
		},
		{
			MethodName: "LoginStudent",
			Handler:    _StudentService_LoginStudent_Handler,
		},
		{
			MethodName: "GetStudentProfile",
			Handler:    _StudentService_GetStudentProfile_Handler,
		},
		{
			MethodName: "UpdateStudentProfile",
			Handler:    _StudentService_UpdateStudentProfile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/education.proto",
}

const (
	EnrollmentService_EnrollStudent_FullMethodName       = "/education.EnrollmentService/EnrollStudent"
	EnrollmentService_GetStudentsByCourse_FullMethodName = "/education.EnrollmentService/GetStudentsByCourse"
	EnrollmentService_GetCoursesByStudent_FullMethodName = "/education.EnrollmentService/GetCoursesByStudent"
	EnrollmentService_UnEnrollStudent_FullMethodName     = "/education.EnrollmentService/UnEnrollStudent"
)

// EnrollmentServiceClient is the client API for EnrollmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnrollmentServiceClient interface {
	EnrollStudent(ctx context.Context, in *EnrollmentRequest, opts ...grpc.CallOption) (*Empty, error)
	GetStudentsByCourse(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*StudentList, error)
	GetCoursesByStudent(ctx context.Context, in *StudentIDRequest, opts ...grpc.CallOption) (*CourseList, error)
	UnEnrollStudent(ctx context.Context, in *EnrollmentRequest, opts ...grpc.CallOption) (*Empty, error)
}

type enrollmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEnrollmentServiceClient(cc grpc.ClientConnInterface) EnrollmentServiceClient {
	return &enrollmentServiceClient{cc}
}

func (c *enrollmentServiceClient) EnrollStudent(ctx context.Context, in *EnrollmentRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, EnrollmentService_EnrollStudent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enrollmentServiceClient) GetStudentsByCourse(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*StudentList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StudentList)
	err := c.cc.Invoke(ctx, EnrollmentService_GetStudentsByCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enrollmentServiceClient) GetCoursesByStudent(ctx context.Context, in *StudentIDRequest, opts ...grpc.CallOption) (*CourseList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseList)
	err := c.cc.Invoke(ctx, EnrollmentService_GetCoursesByStudent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enrollmentServiceClient) UnEnrollStudent(ctx context.Context, in *EnrollmentRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, EnrollmentService_UnEnrollStudent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnrollmentServiceServer is the server API for EnrollmentService service.
// All implementations must embed UnimplementedEnrollmentServiceServer
// for forward compatibility.
type EnrollmentServiceServer interface {
	EnrollStudent(context.Context, *EnrollmentRequest) (*Empty, error)
	GetStudentsByCourse(context.Context, *CourseIDRequest) (*StudentList, error)
	GetCoursesByStudent(context.Context, *StudentIDRequest) (*CourseList, error)
	UnEnrollStudent(context.Context, *EnrollmentRequest) (*Empty, error)
	mustEmbedUnimplementedEnrollmentServiceServer()
}

// UnimplementedEnrollmentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEnrollmentServiceServer struct{}

func (UnimplementedEnrollmentServiceServer) EnrollStudent(context.Context, *EnrollmentRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnrollStudent not implemented")
}
func (UnimplementedEnrollmentServiceServer) GetStudentsByCourse(context.Context, *CourseIDRequest) (*StudentList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentsByCourse not implemented")
}
func (UnimplementedEnrollmentServiceServer) GetCoursesByStudent(context.Context, *StudentIDRequest) (*CourseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCoursesByStudent not implemented")
}
func (UnimplementedEnrollmentServiceServer) UnEnrollStudent(context.Context, *EnrollmentRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnEnrollStudent not implemented")
}
func (UnimplementedEnrollmentServiceServer) mustEmbedUnimplementedEnrollmentServiceServer() {}
func (UnimplementedEnrollmentServiceServer) testEmbeddedByValue()                           {}

// UnsafeEnrollmentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnrollmentServiceServer will
// result in compilation errors.
type UnsafeEnrollmentServiceServer interface {
	mustEmbedUnimplementedEnrollmentServiceServer()
}

func RegisterEnrollmentServiceServer(s grpc.ServiceRegistrar, srv EnrollmentServiceServer) {
	// If the following call pancis, it indicates UnimplementedEnrollmentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EnrollmentService_ServiceDesc, srv)
}

func _EnrollmentService_EnrollStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnrollmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnrollmentServiceServer).EnrollStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EnrollmentService_EnrollStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnrollmentServiceServer).EnrollStudent(ctx, req.(*EnrollmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnrollmentService_GetStudentsByCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnrollmentServiceServer).GetStudentsByCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EnrollmentService_GetStudentsByCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnrollmentServiceServer).GetStudentsByCourse(ctx, req.(*CourseIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnrollmentService_GetCoursesByStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnrollmentServiceServer).GetCoursesByStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EnrollmentService_GetCoursesByStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnrollmentServiceServer).GetCoursesByStudent(ctx, req.(*StudentIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnrollmentService_UnEnrollStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnrollmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnrollmentServiceServer).UnEnrollStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EnrollmentService_UnEnrollStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnrollmentServiceServer).UnEnrollStudent(ctx, req.(*EnrollmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EnrollmentService_ServiceDesc is the grpc.ServiceDesc for EnrollmentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EnrollmentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "education.EnrollmentService",
	HandlerType: (*EnrollmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EnrollStudent",
			Handler:    _EnrollmentService_EnrollStudent_Handler,
		},
		{
			MethodName: "GetStudentsByCourse",
			Handler:    _EnrollmentService_GetStudentsByCourse_Handler,
		},
		{
			MethodName: "GetCoursesByStudent",
			Handler:    _EnrollmentService_GetCoursesByStudent_Handler,
		},
		{
			MethodName: "UnEnrollStudent",
			Handler:    _EnrollmentService_UnEnrollStudent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/education.proto",
}

const (
	LectureService_AddLectureToCourse_FullMethodName  = "/education.LectureService/AddLectureToCourse"
	LectureService_GetLecturesByCourse_FullMethodName = "/education.LectureService/GetLecturesByCourse"
	LectureService_GetLectureContent_FullMethodName   = "/education.LectureService/GetLectureContent"
	LectureService_UpdateLecture_FullMethodName       = "/education.LectureService/UpdateLecture"
	LectureService_DeleteLecture_FullMethodName       = "/education.LectureService/DeleteLecture"
)

// LectureServiceClient is the client API for LectureService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LectureServiceClient interface {
	AddLectureToCourse(ctx context.Context, in *LectureRequest, opts ...grpc.CallOption) (*Lecture, error)
	GetLecturesByCourse(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*LectureList, error)
	GetLectureContent(ctx context.Context, in *LectureIDRequest, opts ...grpc.CallOption) (*LectureContent, error)
	UpdateLecture(ctx context.Context, in *UpdateLectureRequest, opts ...grpc.CallOption) (*Lecture, error)
	DeleteLecture(ctx context.Context, in *LectureIDRequest, opts ...grpc.CallOption) (*Empty, error)
}

type lectureServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLectureServiceClient(cc grpc.ClientConnInterface) LectureServiceClient {
	return &lectureServiceClient{cc}
}

func (c *lectureServiceClient) AddLectureToCourse(ctx context.Context, in *LectureRequest, opts ...grpc.CallOption) (*Lecture, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Lecture)
	err := c.cc.Invoke(ctx, LectureService_AddLectureToCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lectureServiceClient) GetLecturesByCourse(ctx context.Context, in *CourseIDRequest, opts ...grpc.CallOption) (*LectureList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LectureList)
	err := c.cc.Invoke(ctx, LectureService_GetLecturesByCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lectureServiceClient) GetLectureContent(ctx context.Context, in *LectureIDRequest, opts ...grpc.CallOption) (*LectureContent, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LectureContent)
	err := c.cc.Invoke(ctx, LectureService_GetLectureContent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lectureServiceClient) UpdateLecture(ctx context.Context, in *UpdateLectureRequest, opts ...grpc.CallOption) (*Lecture, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Lecture)
	err := c.cc.Invoke(ctx, LectureService_UpdateLecture_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lectureServiceClient) DeleteLecture(ctx context.Context, in *LectureIDRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, LectureService_DeleteLecture_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LectureServiceServer is the server API for LectureService service.
// All implementations must embed UnimplementedLectureServiceServer
// for forward compatibility.
type LectureServiceServer interface {
	AddLectureToCourse(context.Context, *LectureRequest) (*Lecture, error)
	GetLecturesByCourse(context.Context, *CourseIDRequest) (*LectureList, error)
	GetLectureContent(context.Context, *LectureIDRequest) (*LectureContent, error)
	UpdateLecture(context.Context, *UpdateLectureRequest) (*Lecture, error)
	DeleteLecture(context.Context, *LectureIDRequest) (*Empty, error)
	mustEmbedUnimplementedLectureServiceServer()
}

// UnimplementedLectureServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLectureServiceServer struct{}

func (UnimplementedLectureServiceServer) AddLectureToCourse(context.Context, *LectureRequest) (*Lecture, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLectureToCourse not implemented")
}
func (UnimplementedLectureServiceServer) GetLecturesByCourse(context.Context, *CourseIDRequest) (*LectureList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLecturesByCourse not implemented")
}
func (UnimplementedLectureServiceServer) GetLectureContent(context.Context, *LectureIDRequest) (*LectureContent, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLectureContent not implemented")
}
func (UnimplementedLectureServiceServer) UpdateLecture(context.Context, *UpdateLectureRequest) (*Lecture, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLecture not implemented")
}
func (UnimplementedLectureServiceServer) DeleteLecture(context.Context, *LectureIDRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLecture not implemented")
}
func (UnimplementedLectureServiceServer) mustEmbedUnimplementedLectureServiceServer() {}
func (UnimplementedLectureServiceServer) testEmbeddedByValue()                        {}

// UnsafeLectureServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LectureServiceServer will
// result in compilation errors.
type UnsafeLectureServiceServer interface {
	mustEmbedUnimplementedLectureServiceServer()
}

func RegisterLectureServiceServer(s grpc.ServiceRegistrar, srv LectureServiceServer) {
	// If the following call pancis, it indicates UnimplementedLectureServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LectureService_ServiceDesc, srv)
}

func _LectureService_AddLectureToCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LectureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LectureServiceServer).AddLectureToCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LectureService_AddLectureToCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LectureServiceServer).AddLectureToCourse(ctx, req.(*LectureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LectureService_GetLecturesByCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LectureServiceServer).GetLecturesByCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LectureService_GetLecturesByCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LectureServiceServer).GetLecturesByCourse(ctx, req.(*CourseIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LectureService_GetLectureContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LectureIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LectureServiceServer).GetLectureContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LectureService_GetLectureContent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LectureServiceServer).GetLectureContent(ctx, req.(*LectureIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LectureService_UpdateLecture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLectureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LectureServiceServer).UpdateLecture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LectureService_UpdateLecture_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LectureServiceServer).UpdateLecture(ctx, req.(*UpdateLectureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LectureService_DeleteLecture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LectureIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LectureServiceServer).DeleteLecture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LectureService_DeleteLecture_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LectureServiceServer).DeleteLecture(ctx, req.(*LectureIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LectureService_ServiceDesc is the grpc.ServiceDesc for LectureService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LectureService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "education.LectureService",
	HandlerType: (*LectureServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddLectureToCourse",
			Handler:    _LectureService_AddLectureToCourse_Handler,
		},
		{
			MethodName: "GetLecturesByCourse",
			Handler:    _LectureService_GetLecturesByCourse_Handler,
		},
		{
			MethodName: "GetLectureContent",
			Handler:    _LectureService_GetLectureContent_Handler,
		},
		{
			MethodName: "UpdateLecture",
			Handler:    _LectureService_UpdateLecture_Handler,
		},
		{
			MethodName: "DeleteLecture",
			Handler:    _LectureService_DeleteLecture_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/education.proto",
}
