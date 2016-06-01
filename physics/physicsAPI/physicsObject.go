package physicsAPI

import vmath "github.com/walesey/go-engine/vectormath"

type PhysicsObject interface {
	Delete()
	ApplyForce(force, position vmath.Vector3)
	ApplyTorque(torque vmath.Vector3)

	GetPosition() vmath.Vector3
	GetVelocity() vmath.Vector3
	GetOrientation() vmath.Quaternion
	GetAngularVelocityVector() vmath.Vector3
	GetMass() float64
	GetRadius() float64
	GetFriction() float64
	GetRestitution() float64
	InertiaTensor() vmath.Matrix3
	IsStatic() bool

	SetPosition(position vmath.Vector3)
	SetVelocity(velocity vmath.Vector3)
	SetOrientation(orientation vmath.Quaternion)
	SetAngularVelocityVector(av vmath.Vector3)
	SetMass(mass float64)
	SetRadius(radius float64)
	SetFriction(friction float64)
	SetRestitution(restitution float64)
}

type PhysicsObject2D interface {
	ApplyForce(force, position vmath.Vector3)
	ApplyTorque(torque vmath.Vector3)
	GetPosition() vmath.Vector2
	GetVelocity() vmath.Vector2
	GetRotation() float64
	GetAngularVelocity() float64
	GetMass() float64
	SetVelocity(velocity vmath.Vector2)
	SetRotation(orientation vmath.Quaternion)
	SetAngularVelocity(av float64)
}

type PhysicsShape2D interface {
}
